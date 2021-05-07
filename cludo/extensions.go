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
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable CLUDO_NUMBER. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(config.EnvPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viper.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		logDebugf("Flag Name: %v", f.Name)
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --server-url to CLUDO_SERVER_URL
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			logDebugf("ENV VAR: %s_%s\n", config.EnvPrefix, envVarSuffix)
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", config.EnvPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			//path := "client.default." + strings.ReplaceAll(f.Name, "-", "_")
			logDebugf("Flag Name (%s) & New Value (%s)\n", f.Name, val)
			//fmt.Printf("%#v", cmd.Flags())
			// This doesn't appear to be effective
			// Maybe something with the way we are accessing the data?
			// the flag path is 'key-path' versus 'config.default.key_path'
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))

			//viper.Set(path, fmt.Sprintf("%v", val))
		}
	})
}
