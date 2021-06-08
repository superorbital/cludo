package cludo

import (
	"github.com/spf13/cobra"
	"github.com/superorbital/cludo/pkg/config"
)

// MakeExecCmd sets up the exec subcommand.
func MakeExecCmd(debug bool, dryRun bool, exit func(int)) (*cobra.Command, error) {
	var cleanEnv bool

	execCmd := &cobra.Command{
		Use:   "exec",
		Short: "Executes a command with an environment setup with cludo credentials",
		Long:  `Executes a command with an environment setup with cludo credentials for the configured cludo client.`,
		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := config.NewConfigFromViper()
			cobra.CheckErr(err)
			clientConfig := userConfig.Client

			bundle, err := GenerateEnvironment(clientConfig, userConfig.Target, debug, dryRun)
			cobra.CheckErr(err)

			code, err := ExecWithEnv(args, bundle, !cleanEnv, userConfig.Target)
			cobra.CheckErr(err)

			if code != 0 {
				exit(code)
			}
		},
	}
	execCmd.Flags().BoolVar(&cleanEnv, "clean-env", false, "Set to run shell without inheriting from current environment")

	return execCmd, nil
}
