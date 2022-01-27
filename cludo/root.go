package cludo

import (
	"fmt"

	"github.com/superorbital/cludo/pkg/build"
	"github.com/superorbital/cludo/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MakeRootCmd sets up the root cmd and all subcommands.
func MakeRootCmd(exit func(int)) (*cobra.Command, error) {
	var configFile string
	var debug bool
	var dryRun bool

	// Use executable name as the command name
	rootCmd := &cobra.Command{
		Version: build.VersionFull(),
		Use:     config.CludoExecutable,
		Short:   "Cloud Sudo Client CLI",
		Long:    `This is the Cloud Sudo Client CLI, which end users will typically use to interact with the Cloud Sudo server.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return config.ConfigureViper(config.CludoExecutable, configFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Register flags bound to viper.
	rootCmd.PersistentFlags().String("target", "", "URL of server appended with the role name")
	viper.BindPFlag("target", rootCmd.PersistentFlags().Lookup("target"))
	rootCmd.PersistentFlags().Bool("interactive", true, "Should the CLI prompt the user for input")
	viper.BindPFlag("client.default.interactive", rootCmd.PersistentFlags().Lookup("interactive"))
	rootCmd.PersistentFlags().StringSliceP("ssh-key-paths", "k", []string{"~/.ssh/id_ed25519"}, "A comma seperated list of SSH private key filesystem paths")
	viper.BindPFlag("ssh_key_paths", rootCmd.PersistentFlags().Lookup("ssh-key-paths"))

	// Register regular flags.
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to cludo.yaml config file. Overrides normal search paths.")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enables debug logging")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Prints requests to make to server instead of actually sending them")

	// cludo env
	envCmd, err := MakeEnvCmd(debug, dryRun)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up env: %v", err)
	}
	rootCmd.AddCommand(envCmd)

	// cludo exec
	execCmd, err := MakeExecCmd(debug, dryRun, exit)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up exec: %v", err)
	}
	rootCmd.AddCommand(execCmd)

	// cludo shell
	shellCmd, err := MakeShellCmd(debug, dryRun, exit)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up shell: %v", err)
	}
	rootCmd.AddCommand(shellCmd)

	// cludo version
	versionCmd, err := MakeVersionCmd(debug, dryRun)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up version: %v", err)
	}
	rootCmd.AddCommand(versionCmd)

	// cludo lint
	lintCmd, err := MakeLintCmd(debug, dryRun)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up lint")
	}
	rootCmd.AddCommand(lintCmd)

	// cludo completion [...]
	completionCmd := MakeCompletionCmd()
	rootCmd.AddCommand(completionCmd)

	return rootCmd, nil
}
