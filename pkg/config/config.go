package config

type Config struct {
	Client      *ClientConfig `mapstructure:"client"`
	Server      *ServerConfig `mapstructure:"server"`
	Target      string        `mapstructure:"target"`
	SSHKeyPaths []string      `mapstructure:"ssh_key_paths"`
}
