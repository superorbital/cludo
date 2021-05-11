package cludo

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/superorbital/cludo/pkg/config"
)

func initializeConfig(cmd *cobra.Command) error {
	initViperConfigs()

	err := bindEnvVars(cmd)
	if err != nil {
		return err
	}
	return nil
}

func bindEnvVars(cmd *cobra.Command) error {

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --profile
	// binds to an environment variable CLUDO_PROFILE. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(config.EnvPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --server-url which we fix in the bindFlags function
	viper.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		//logDebugf("Flag Name: %v", f.Name)
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --server-url to CLUDO_SERVER_URL
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			//logDebugf("ENV VAR: %s_%s\n", config.EnvPrefix, envVarSuffix)
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", config.EnvPrefix, envVarSuffix))
		}

		// FIXME: This is a hacked together mess at the moment.
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			//path := "client.default." + strings.ReplaceAll(f.Name, "-", "_")
			logDebugf("Flag Name (%s) & New Value (%s)\n", f.Name, val)
			//fmt.Printf("%#v", cmd.Flags())
			// How are we supposed to use this?
			// This is the intended action, but it is not working as expected.
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))

			// This is the current workaround.
			switch f.Name {
			case "server-url":
				userConfig.Client["default"].ServerURL = fmt.Sprintf("%v", val)
			case "key-path":
				userConfig.Client["default"].KeyPath = fmt.Sprintf("%v", val)
			case "shell-path":
				userConfig.Client["default"].ShellPath = fmt.Sprintf("%v", val)
			case "passphrase":
				userConfig.Client["default"].Passphrase = fmt.Sprintf("%v", val)
			}
		}
	})
}
