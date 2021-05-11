package cludo

import (
	"fmt"

	"github.com/superorbital/cludo/pkg/build"
	"github.com/superorbital/cludo/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MakeRootCmd sets up the root cmd and all subcommands.
func MakeRootCmd() (*cobra.Command, error) {
	var configFile string
	var debug bool
	var dryRun bool
	var profile string

	// Use executable name as the command name
	rootCmd := &cobra.Command{
		Version: build.Version,
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
	rootCmd.Flags().StringP("key-path", "k", "~/.ssh/id_rsa", "Path to SSH private key")
	viper.BindPFlag("client.default.key_path", rootCmd.Flags().Lookup("key-path"))
	rootCmd.Flags().StringP("passphrase", "a", "", "Passphrase for private key (consider using config or setting CLUDO_PASSPHRASE)")
	viper.BindPFlag("client.default.passphrase", rootCmd.Flags().Lookup("passphrase"))

	// Register regular flags.
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to cludo.yaml config file. Overrides normal search paths.")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enables debug logging")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Prints requests to make to server instead of actually sending them")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Connection profile name")

	// Add all subcommands
	envCmd, err := MakeEnvCmd(debug, dryRun, profile)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up exec")
	}
	rootCmd.AddCommand(envCmd)

	rootCmd.AddCommand(MakeCompletionCmd())

	return rootCmd, nil
}
