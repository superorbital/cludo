package config

import "github.com/spf13/viper"

const DefaultClientConfig = "default"

type Config struct {
	Client map[string]*ClientConfig `yaml:"client"`
	Server *ServerConfig            `yaml:"server"`
}

func NewConfigFromViper() (*Config, error) {
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
