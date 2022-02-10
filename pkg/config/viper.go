package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const CludodExecutable = "cludod"
const CludoExecutable = "cludo"

// envPrefix defines a prefix required on all environment variables bound to Config.
// For example the path "number" is bound to CLUDO_NUMBER.
const envPrefix = "CLUDO"

var errConfigNotFound = errors.New("Failed to load configuration file: File not found")

func errConfigLoadFailed(cause error) error {
	return fmt.Errorf("Failed to load configuration file: %v", cause)
}

func errHomeDirFailed(cause error) error {
	return fmt.Errorf("Failed to get user home directory: %v", cause)
}

func errWorkingDirFailed(cause error) error {
	return fmt.Errorf("Failed to get working directory: %v", cause)
}

func ConfigureViper(executable string, configFile string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errWorkingDirFailed(err)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return errHomeDirFailed(err)
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --profile
	// binds to an environment variable CLUDO_PROFILE. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	// Read configuration
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(executable)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path.Join("/etc", executable))
		viper.AddConfigPath(path.Join(homedir, fmt.Sprintf(".%s", executable)))
		viper.AddConfigPath(path.Join(homedir, ".config", executable))
		viper.AddConfigPath(cwd)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return errConfigNotFound
		}
		return errConfigLoadFailed(err)
	}
	// local repo cludo.yaml file
	// Only check for this when we are using the client.
	if executable == "cludo" {
		viper.AddConfigPath(".")
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return errConfigNotFound
			}
			return errConfigLoadFailed(err)
		}
	}
	return nil
}

func NewConfigFromViper() (*Config, error) {
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
