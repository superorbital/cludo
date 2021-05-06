package config

import "github.com/spf13/viper"

const (
	CludoVersion        = "v0.0.1"
	DefaultClientConfig = "default"

	// The name of our config file, without the file extension because viper supports many different config file languages.
	DefaultConfigFilename = "cludo"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to CLUDO_NUMBER.
	EnvPrefix = "CLUDO"
)

type Config struct {
	Client map[string]*ClientConfig `mapstructure:"client"`
	Server *ServerConfig            `mapstructure:"server"`
}

func NewConfigFromViper() (*Config, error) {
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
