package config

type Config struct {
	Client               map[string]*ClientConfig `mapstructure:"client"`
	ClientProfileDefault string                   `mapstructure:"client_profile_default"`
	Server               *ServerConfig            `mapstructure:"server"`
}
