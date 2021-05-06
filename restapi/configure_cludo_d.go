// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
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

var cludoDFlags = struct {
	Example1 string `long:"example1" description:"Sample for showing how to configure cmd-line flags"`
	Example2 string `long:"example2" description:"Further info at https://github.com/jessevdk/go-flags"`
}{}

func configureFlags(api *operations.CludoDAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "cludod Flags",
			LongDescription:  "",
			Options:          &cludoDFlags,
		},
	}
}

func configureAPI(api *operations.CludoDAPI) http.Handler {
	// Read configuration
	viper.SetConfigName("cludo")         // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/cludod/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.cludod") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			api.Logger("ERROR: Failed to load configuration file: File not found")
			return nil
		} else {
			// Config file was found but another error was produced
			api.Logger("ERROR: Failed to load configuration file:", err)
			return nil
		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		api.Logger("Config file changed:", e.Name)
	})

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "X-CLUDO-KEY" header is set
	if api.APIKeyHeaderAuth == nil {
		api.APIKeyHeaderAuth = func(token string) (*models.ModelsPrincipal, error) {
			conf, err := config.NewConfigFromViper()
			if err != nil {
				return nil, errors.New(500, "Failed to read cludo configuration: %v", err)
			}

			if conf.Server == nil {
				return nil, errors.New(500, "Server configuration is missing")
			}

			authz, err := conf.Server.NewAuthorizer()
			if err != nil {
				return nil, errors.New(500, "Failed to create authorizer: %v", err)
			}

			id, ok, err := authz.CheckAuthHeader(token)
			if err != nil {
				return nil, errors.New(500, "Failed to validate message signature: %v", err)
			}
			if ok {
				for _, user := range conf.Server.Users {
					if user.ID() == id {
						roles := models.ModelsPrincipalRoles{
							Aws: models.ModelsAWSPrincipalRoles{},
						}

						for roleID, role := range user.Roles.AWS {
							roles.Aws[roleID] = models.ModelsAWSPrincipalRole{
								SessionDuration: role.SessionDuration.String(),
								AccessKeyID:     role.AccessKeyID,
								SecretAccessKey: role.SecretAccessKey,
								Arn:             role.AssumeRoleARN,
							}
						}

						return &models.ModelsPrincipal{
							PublicKey:   user.PublicKey,
							DefaultRole: user.DefaultRole,
							Roles:       &roles,
						}, nil
					}
				}
			}
			return nil, errors.Unauthenticated("APIKeyHeaderAuth")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.EnvironmentGenerateEnvironmentHandler == nil {
		api.EnvironmentGenerateEnvironmentHandler = environment.GenerateEnvironmentHandlerFunc(func(params environment.GenerateEnvironmentParams, principal *models.ModelsPrincipal) middleware.Responder {
			config := config.Config{}
			err := viper.Unmarshal(&config)
			if err != nil {
				api.Logger("ERROR: Failed to read config:", err)
				return middleware.Error(500, &struct{}{})
			}

			return middleware.NotImplemented("operation environment.GenerateEnvironment has not yet been implemented")
		})
	}
	if api.SystemHealthHandler == nil {
		api.SystemHealthHandler = system.HealthHandlerFunc(func(params system.HealthParams) middleware.Responder {
			return system.NewHealthOK().WithPayload(&models.ModelsHealthResponse{
				Status:  true,
				Version: build.Version,
			})
		})
	}
	if api.RoleListRolesHandler == nil {
		api.RoleListRolesHandler = role.ListRolesHandlerFunc(func(params role.ListRolesParams, principal *models.ModelsPrincipal) middleware.Responder {
			conf, err := config.NewConfigFromViper()
			if err != nil {
				errMsg := fmt.Sprintf("Failed to read cludo configuration: %v", err)
				return role.NewListRolesDefault(500).WithPayload(&models.Error{
					Code:    500,
					Message: &errMsg,
				})
			}

			if conf.Server == nil {
				errMsg := fmt.Sprintf("Server configuration is missing")
				return role.NewListRolesDefault(500).WithPayload(&models.Error{
					Code:    500,
					Message: &errMsg,
				})
			}

			roles := []string{}
			if principal != nil && principal.Roles != nil && principal.Roles.Aws != nil {
				for roleID, _ := range principal.Roles.Aws {
					roles = append(roles, roleID)
				}
			}

			return role.NewListRolesOK().WithPayload(&models.ModelsRoleIDsResponse{
				Roles: roles,
			})
		})
	}

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
	return handler
}
