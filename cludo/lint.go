package cludo

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/superorbital/cludo/pkg/config"
)

// MakeLintCmd sets up the lint subcommand.
func MakeLintCmd(debug bool, dryRun bool, profile string) (*cobra.Command, error) {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Check cludo's configuration's validity",
		Long:  `Check cludo configuration's validity. Also prints current values.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			_, err := config.NewConfigFromViper()
			cobra.CheckErr(err)

			_, err = fmt.Fprintln(out, PrintConfig())
			cobra.CheckErr(err)
		},
	}

	return lintCmd, nil
}

// PrintConfig turns a bundle into a shell script that can be sourced to export environment variables.
func PrintConfig() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Cludo config:\n\n")
	for _, v := range viper.AllKeys() {
		fmt.Fprintf(&b, "%s=%#v\n", v, viper.Get(v))
	}
	return b.String()
}
