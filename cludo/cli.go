package cludo

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/superorbital/cludo/client"
	"github.com/superorbital/cludo/pkg/config"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// debug flag indicating that cli should output debug logs
var debug bool

// config file location
var configFile string

// dry run flag
var dryRun bool

// name of the executable
var exeName string = filepath.Base(os.Args[0])

// Config
var clientProfile string
var userConfig config.Config

// logDebugf writes debug log to stdout
func logDebugf(format string, v ...interface{}) {
	if !debug {
		return
	}
	log.Printf(format, v...)
}

// depth of recursion to construct model flags
var maxDepth int = 5

// makeClient constructs a client object
func makeClient(cmd *cobra.Command, args []string) (*client.Cludod, error) {
	hostname := viper.GetString("hostname")
	scheme := viper.GetString("scheme")

	r := httptransport.New(hostname, client.DefaultBasePath, []string{scheme})
	r.SetDebug(debug)
	// set custom producer and consumer to use the default ones

	r.Consumers["application/json"] = runtime.JSONConsumer()

	r.Producers["application/json"] = runtime.JSONProducer()

	appCli := client.New(r, strfmt.Default)
	logDebugf("Server url: %v://%v", scheme, hostname)
	return appCli, nil
}

// MakeRootCmd returns the root cmd
func MakeRootCmd() (*cobra.Command, error) {
	// Note: The first time this is called debug will still be false, so no Debug logs.
	// Might consider an ENV Var check to make this immediate.
	cobra.OnInitialize(initViperConfigs)

	// Use executable name as the command name
	rootCmd := &cobra.Command{
		Version: config.CludoVersion,
		Use:     exeName,
		Short:   "Cloud Sudo Client CLI",
		Long:    `This is the Cloud Sudo Client CLI, which end users will typically use to interact with the Cloud Sudo server.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// register basic flags
	rootCmd.PersistentFlags().String("hostname", client.DefaultHost, "hostname of the service")
	viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname"))
	rootCmd.PersistentFlags().String("scheme", client.DefaultSchemes[0], fmt.Sprintf("Choose from: %v", client.DefaultSchemes))
	viper.BindPFlag("scheme", rootCmd.PersistentFlags().Lookup("scheme"))

	// configure debug flag
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "output debug logs")
	// configure config location
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file path")
	// configure dry run flag
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "do not send the request to server")

	// register security flags
	// add all operation groups
	operationGroupExecCmd, err := makeOperationGroupExecCmd()
	if err != nil {
		return nil, err
	}

	rootCmd.AddCommand(operationGroupExecCmd)
	operationGroupExecCmd.Flags().StringVarP(&clientProfile, "profile", "p", "default", "Connection profile name")
	operationGroupExecCmd.Flags().StringP("server-url", "s", "http://localhost:80/", "Complete server URL")
	viper.BindPFlag("client.default.server_url", operationGroupExecCmd.Flags().Lookup("server-url"))
	operationGroupExecCmd.Flags().StringP("key-path", "k", "~/.ssh/id_rsa", "Path to SSH private key")
	viper.BindPFlag("client.default.key_path", operationGroupExecCmd.Flags().Lookup("key-path"))
	operationGroupExecCmd.Flags().StringP("passphrase", "a", "", "Passphrase for private key (consider using config or setting CLUDO_PASSPHRASE)")
	viper.BindPFlag("client.default.passphrase", operationGroupExecCmd.Flags().Lookup("passphrase"))
	operationGroupExecCmd.Flags().StringP("shell-path", "i", "/bin/sh", "Path to shell")
	viper.BindPFlag("client.default.shell_path", operationGroupExecCmd.Flags().Lookup("shell-path"))
	operationGroupExecCmd.Flags().StringArrayP("roles", "r", []string{"default"}, "One or more comma seperated roles")
	viper.BindPFlag("client.default.roles", operationGroupExecCmd.Flags().Lookup("roles"))

	// add cobra completion
	rootCmd.AddCommand(makeGenCompletionCmd())

	err = bindEnvVars(rootCmd)
	if err != nil {
		return nil, err
	}

	return rootCmd, nil
}

// initViperConfigs initialize viper config using config file in '$HOME/.config/<cli name>/config.<json|yaml...>'
// currently hostname, scheme and auth tokens can be specified in this config file.
func initViperConfigs() {
	if configFile != "" {
		// use user specified config file location
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("cludo") // name of config file (without extension)
		viper.SetConfigType("yaml")  // REQUIRED if the config file does not have the extension in the name

		// look for default config
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		// Find current working directory.
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("/etc/cludo")
		viper.AddConfigPath(path.Join(home, ".cludo"))
		viper.AddConfigPath(path.Join(home, ".config", exeName)) // go-swagger's default location
		viper.AddConfigPath(path.Join(cwd, ".cludo"))
		viper.AddConfigPath(cwd)
		viper.SetConfigName(config.DefaultConfigFilename)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("[WARN] Could not load config file: %v\n", err)
			return
		}
		logDebugf("Using config file: %v", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&userConfig)
	logDebugf("Client settings from config file: %v", *userConfig.Client["default"])
	logDebugf("Viper Keys: %v", viper.AllKeys())
}

func makeOperationGroupExecCmd() (*cobra.Command, error) {
	operationGroupExecCmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute command with authorization",
		Long:  `Execute a single command with authentication env vars set.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			msg := runOperationExec(cmd, args, userConfig)

			fmt.Fprintln(out, msg)

		},
	}

	//operationWithProfileCmd, err := makeOperationExecWithProfileCmd()
	//if err != nil {
	//	return nil, err
	//}
	//operationGroupExecCmd.AddCommand(operationWithProfileCmd)

	return operationGroupExecCmd, nil
}
