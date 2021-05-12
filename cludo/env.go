package cludo

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/superorbital/cludo/client/environment"
	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/config"
)

// MakeEnvCmd sets up the env subcommand.
func MakeEnvCmd(debug bool, dryRun bool, profile string) (*cobra.Command, error) {
	execCmd := &cobra.Command{
		Use:   "env",
		Short: "Get environment variables for cludo",
		Long:  `Get environment variables for the provided cludo profile (or 'default'). You can add these to your shell with: . $(cludo env)`,
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			userConfig, err := config.NewConfigFromViper()
			cobra.CheckErr(err)

			bundle, err := GenerateEnvironment(userConfig.Client[profile], debug, dryRun)
			cobra.CheckErr(err)

			_, err = fmt.Fprintln(out, FormatBundle(bundle))
			cobra.CheckErr(err)
		},
	}
	execCmd.Flags().StringP("shell-path", "i", "/bin/sh", "Path to shell")
	viper.BindPFlag("client.default.shell_path", execCmd.Flags().Lookup("shell-path"))
	execCmd.Flags().StringArrayP("roles", "r", []string{"default"}, "One or more comma seperated roles")
	viper.BindPFlag("client.default.roles", execCmd.Flags().Lookup("roles"))

	return execCmd, nil
}

// GenerateEnvironment generates an environment bundle on a remote cludod service.
func GenerateEnvironment(cc *config.ClientConfig, debug bool, dryRun bool) (models.ModelsEnvironmentBundle, error) {
	cludodClient, err := cc.NewClient(debug)
	if err != nil {
		return nil, err
	}
	signer, err := cc.NewDefaultSigner()
	if err != nil {
		return nil, err
	}

	params := environment.NewGenerateEnvironmentParams().WithDefaults().WithBody(&models.ModelsEnvironmentRequest{})

	if dryRun {
		log.Printf("Dry run enabled: Would send %#v", params)
		return nil, nil
	}

	response, err := cludodClient.Environment.GenerateEnvironment(params, signer.CludoAuth())
	if err != nil {
		return nil, fmt.Errorf("[1] Failed to generate environment: %v, %#v", err, response)
	}

	if response != nil && response.Payload != nil {
		return response.Payload.Bundle, nil
	}

	return nil, nil
}

// FormatBundle turns a bundle into a shell script that can be sourced to export environment variables.
func FormatBundle(bundle models.ModelsEnvironmentBundle) string {
	var b strings.Builder
	for k, v := range bundle {
		escaped := strings.ReplaceAll(strings.ReplaceAll(v, "\\", "\\\\"), "\"", "\\\"")
		fmt.Fprintf(&b, "export %s=\"%s\"\n", k, escaped)
	}
	return b.String()
}
