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
	var profile string

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
	rootCmd.PersistentFlags().String("server-url", "", "Base URL of the cludod service to interact with")
	viper.BindPFlag("client.default.server_url", rootCmd.PersistentFlags().Lookup("server-url"))
	rootCmd.PersistentFlags().StringP("key-path", "k", "~/.ssh/id_rsa", "Path to SSH private key")
	viper.BindPFlag("client.default.key_path", rootCmd.PersistentFlags().Lookup("key-path"))
	rootCmd.PersistentFlags().StringP("passphrase", "a", "", "Passphrase for private key (consider using config or setting CLUDO_PASSPHRASE)")
	viper.BindPFlag("client.default.passphrase", rootCmd.PersistentFlags().Lookup("passphrase"))
	rootCmd.PersistentFlags().StringArrayP("roles", "r", []string{"default"}, "One or more comma seperated roles")
	viper.BindPFlag("client.default.roles", rootCmd.PersistentFlags().Lookup("roles"))

	// Register regular flags.
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to cludo.yaml config file. Overrides normal search paths.")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enables debug logging")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Prints requests to make to server instead of actually sending them")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Connection profile name")

	// cludo env
	envCmd, err := MakeEnvCmd(debug, dryRun, profile)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up env: %v", err)
	}
	rootCmd.AddCommand(envCmd)

	// cludo exec
	execCmd, err := MakeExecCmd(debug, dryRun, profile, exit)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up exec: %v", err)
	}
	rootCmd.AddCommand(execCmd)

	// cludo shell
	shellCmd, err := MakeShellCmd(debug, dryRun, profile, exit)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up shell: %v", err)
	}
	rootCmd.AddCommand(shellCmd)

	// cludo version
	versionCmd, err := MakeVersionCmd(debug, dryRun, profile)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up version: %v", err)
	}
	rootCmd.AddCommand(versionCmd)

	// cludo lint
	lintCmd, err := MakeLintCmd(debug, dryRun, profile)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up lint")
	}
	rootCmd.AddCommand(lintCmd)

	// cludo completion [...]
	completionCmd := MakeCompletionCmd()
	rootCmd.AddCommand(completionCmd)

	return rootCmd, nil
}
