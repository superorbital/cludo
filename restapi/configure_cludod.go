// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"

	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/build"
	"github.com/superorbital/cludo/pkg/config"
	"github.com/superorbital/cludo/restapi/operations"
	"github.com/superorbital/cludo/restapi/operations/environment"
	"github.com/superorbital/cludo/restapi/operations/role"
	"github.com/superorbital/cludo/restapi/operations/system"
)

//go:generate swagger generate server --target ../../cludo --name CludoD --spec ../swagger.yaml --principal interface{}

var CludodFlags = struct {
	Config string `long:"config" description:"Path to a configuration file to load."`
}{}

func configureFlags(api *operations.CludodAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "cludod Flags",
			LongDescription:  "",
			Options:          &CludodFlags,
		},
	}
}

func configureAPI(api *operations.CludodAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = log.Printf

	config.ConfigureViper(config.CludodExecutable, CludodFlags.Config)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		api.Logger("[WARN] Config file changed:", e.Name)
	})

	// api.UseSwaggerUI()
	api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "X-CLUDO-KEY" header is set
	api.APIKeyHeaderAuth = func(token string) (*models.ModelsPrincipal, error) {
		conf, err := config.NewConfigFromViper()
		if err != nil {
			api.Logger("[ERROR] Failed to read cludod configuration: %v", err)
			return nil, errors.New(500, "Failed to read cludod configuration: %v", err)
		}

		if conf.Server == nil {
			api.Logger("[ERROR] Server configuration is missing")
			return nil, errors.New(500, "Server configuration is missing")
		}

		authz, err := conf.Server.NewConfigAuthorizer()
		if err != nil {
			api.Logger("[ERROR] Failed to create internal authorizer from server config: %v", err)
			return nil, errors.New(500, "Failed to create internal authorizer from server config: %v", err)
		}

		id, ok, err := authz.CheckAuthHeader(token)
		if err != nil {
			api.Logger("[ERROR] Failed to validate API request signature in header: %v", err)
			return nil, errors.New(500, "Failed to validate API request signature in header: %v", err)
		}
		if ok {
			for _, user := range conf.Server.Users {
				if user.ID() == id {
					principalID := models.ModelsPrincipal(id)
					return &principalID, nil
				}
			}
		} else {
			authz, err := conf.Server.NewGithubAuthorizer()
			if err != nil {
				api.Logger("[ERROR] Failed to create internal authorizer from Github API: %v", err)
				return nil, errors.New(500, "Failed to create internal authorizer from Github API: %v", err)
			}
			id, ok, err := authz.CheckAuthHeader(token)
			if err != nil {
				api.Logger("[ERROR] Failed to validate API request signature in header against Github: %v", err)
				return nil, errors.New(500, "Failed to validate API request signature in header against Github: %v", err)
			}
			if ok {
				for _, user := range conf.Server.Users {
					if user.ID() == id {
						principalID := models.ModelsPrincipal(id)
						return &principalID, nil
					}
				}
			}
		}
		return nil, errors.Unauthenticated("APIKeyHeaderAuth")
	}

	api.EnvironmentGenerateEnvironmentHandler = environment.GenerateEnvironmentHandlerFunc(func(params environment.GenerateEnvironmentParams, principal *models.ModelsPrincipal) middleware.Responder {
		conf, err := config.NewConfigFromViper()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to read cludod configuration: %v", err)
			api.Logger("[ERROR] %s", err)
			return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}

		if conf.Server == nil {
			errMsg := fmt.Sprintf("Server configuration is missing")
			api.Logger("[ERROR] %s", err)
			return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}

		var role *config.AWSRoleConfig
		user, ok := conf.Server.GetUser(string(*principal))
		requestedTargetURI, err := url.Parse(params.Body.Target)
		name := user.Name
		if name == "" {
			name = "UNKNOWN"
		}
		trunc_pubkey := user.PublicKey[0:40] + "..." + user.PublicKey[len(user.PublicKey)-40:]
		remoteIP := GetIP(params.HTTPRequest)
		api.Logger("[AUDIT] Someone from (%s) matching %s's pubkey {%s} is attempting to authenticate to: [%s]", remoteIP, name, trunc_pubkey, requestedTargetURI)
		if err != nil {
			errMsg := fmt.Sprintf("Expected target in URL format, received: %s", params.Body.Target)
			api.Logger("[ERROR] %s", err)
			return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}
		target := strings.TrimLeft(requestedTargetURI.Path, "/")
		if ok && user != nil {
			validTarget := false
			for _, validUserTarget := range user.Targets {
				if target == validUserTarget {
					validTarget = true
				}
			}
			if !validTarget {
				errMsg := fmt.Sprintf("User does not have access to requested target: %s", target)
				api.Logger("[ERROR] %s", errMsg)
				return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
					Code:    403,
					Message: &errMsg,
				})
			}
			role = conf.Server.Targets[target].AWS
		}

		if role == nil {
			errMsg := fmt.Sprintf("Failed to find any roles for user: %v", *principal)
			api.Logger("[ERROR] %s", err)
			return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}

		ap, err := role.NewPlugin()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to initialize plugin: %v", err)
			api.Logger("[ERROR] %s", err)
			return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}
		payload, err := ap.GenerateEnvironment()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to generate environment: %v", err)
			api.Logger("[ERROR] %s", err)
			return environment.NewGenerateEnvironmentDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}

		api.Logger("[AUDIT] Someone from (%s) matching %s's pubkey '%s' authenticated to target: [%s]", remoteIP, name, trunc_pubkey, target)

		return environment.NewGenerateEnvironmentOK().WithPayload(payload)
	})
	api.SystemHealthHandler = system.HealthHandlerFunc(func(params system.HealthParams) middleware.Responder {
		return system.NewHealthOK().WithPayload(&models.ModelsHealthResponse{
			Status:  true,
			Version: build.VersionFull(),
		})
	})
	api.RoleListRolesHandler = role.ListRolesHandlerFunc(func(params role.ListRolesParams, principal *models.ModelsPrincipal) middleware.Responder {
		conf, err := config.NewConfigFromViper()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to read cludod configuration: %v", err)
			api.Logger("[ERROR] %s", err)
			return role.NewListRolesDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}

		if conf.Server == nil {
			errMsg := "Server configuration is missing"
			api.Logger("[ERROR] %s", err)
			return role.NewListRolesDefault(500).WithPayload(&models.Error{
				Code:    500,
				Message: &errMsg,
			})
		}

		roles := []string{}
		user, ok := conf.Server.GetUser(string(*principal))
		if ok && user != nil {
			for _, userTarget := range user.Targets {
				roles = append(roles, conf.Server.Targets[userTarget].AWS.AssumeRoleARN)
			}
		}

		return role.NewListRolesOK().WithPayload(&models.ModelsRoleIDsResponse{
			Roles: roles,
		})
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, handler)
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
