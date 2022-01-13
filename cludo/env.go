package cludo

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/config"
)

// MakeEnvCmd sets up the env subcommand.
func MakeEnvCmd(debug bool, dryRun bool) (*cobra.Command, error) {
	envCmd := &cobra.Command{
		Use:   "env",
		Short: "Get environment variables for cludo",
		Long:  `Get environment variables for the configured cludo target. You can add these to your shell with: . $(cludo env)`,
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			userConfig, err := config.NewConfigFromViper()
			cobra.CheckErr(err)
			bundle, err := GenerateEnvironment(userConfig.Client, userConfig.Target, userConfig.SSHKeyPaths, debug, dryRun)
			cobra.CheckErr(err)

			_, err = fmt.Fprintln(out, FormatBundle(bundle))
			cobra.CheckErr(err)
		},
	}

	return envCmd, nil
}

// FormatBundle turns a bundle into a shell script that can be sourced to export environment variables.
func FormatBundle(bundle models.ModelsEnvironmentBundle) string {
	var b strings.Builder
	for k, v := range bundle {
		if v != "" {
			escaped := strings.ReplaceAll(strings.ReplaceAll(v, "\\", "\\\\"), "\"", "\\\"")
			fmt.Fprintf(&b, "export %s=\"%s\"\n", k, escaped)
		}
	}
	return b.String()
}
