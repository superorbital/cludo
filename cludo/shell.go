package cludo

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/superorbital/cludo/pkg/config"
)

// MakeShellCmd sets up the shell subcommand.
func MakeShellCmd(debug bool, dryRun bool, profile string, exit func(int)) (*cobra.Command, error) {
	var cleanEnv bool

	execCmd := &cobra.Command{
		Use:   "shell",
		Short: "Executes a command with an environment setup with cludo credentials",
		Long:  `Executes a command with an environment setup with cludo credentials for the provided cludo profile (or 'default').`,
		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := config.NewConfigFromViper()
			cobra.CheckErr(err)
			clientConfig := userConfig.Client[profile]

			bundle, err := GenerateEnvironment(clientConfig, userConfig.Target, debug, dryRun)
			cobra.CheckErr(err)

			shell, err := GetShellPath(clientConfig)
			cobra.CheckErr(err)

			code, err := ExecWithEnv(append([]string{shell}, args...), bundle, !cleanEnv, profile, clientConfig.ServerURL)
			cobra.CheckErr(err)

			if code != 0 {
				exit(code)
			}
		},
	}
	execCmd.Flags().BoolVar(&cleanEnv, "clean-env", false, "Set to run shell without inheriting from current environment")
	execCmd.PersistentFlags().StringP("shell-path", "i", "/bin/sh", "Path to shell")
	viper.BindPFlag("client.default.shell_path", execCmd.PersistentFlags().Lookup("shell-path"))

	return execCmd, nil
}

var ErrShellUndefined = errors.New("SHELL undefined!")

func GetShellPath(cc *config.ClientConfig) (string, error) {
	if cc.ShellPath != "" {
		return cc.ShellPath, nil
	}

	envShell := os.Getenv("SHELL")
	if envShell != "" {
		return envShell, nil
	}

	return "", ErrShellUndefined
}
