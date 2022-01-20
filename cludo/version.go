package cludo

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superorbital/cludo/client/system"
	"github.com/superorbital/cludo/pkg/build"
	"github.com/superorbital/cludo/pkg/config"
)

// MakeVersionCmd sets up the exec subcommand.
func MakeVersionCmd(debug bool, dryRun bool) (*cobra.Command, error) {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Get cludo client and cludod server version",
		Long:  `Get cludo client and cludod server version. Uses the server health check endpoint.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			userConfig, err := config.NewConfigFromViper()
			CheckErr(err)

			serverVersion, err := GetVersion(userConfig.Target, debug, dryRun)
			CheckErr(err)

			_, err = fmt.Fprintf(out, "Client: %s\nServer: %s\n", build.VersionFull(), serverVersion)
			CheckErr(err)
		},
	}

	return versionCmd, nil
}

// GetVersion generates an environment bundle on a remote cludod service.
func GetVersion(target string, debug bool, dryRun bool) (string, error) {
	cludodClient, err := config.NewClient(target, debug)
	if err != nil {
		return "", err
	}

	params := system.NewHealthParams()

	if dryRun {
		fmt.Printf("[INFO] Dry run enabled: Would send %#v", params)
		return "", nil
	}

	response, err := cludodClient.System.Health(params)
	if err != nil {
		return "", fmt.Errorf("Failed to check system health: %v, %#v", err, response)
	}

	if response != nil && response.Payload != nil {
		return response.Payload.Version, nil
	}

	return "", nil
}
