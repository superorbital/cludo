package config

import "time"

const DefaultClientConfig = "default"

type AWSRoleConfig struct {
	SessionDuration time.Duration `yaml:"session_duration"`
	AccessKeyID     string        `yaml:"access_key_id"`
	SecretAccessKey string        `yaml:"secret_access_key"`
	AssumeRoleARN   string        `yaml:"arn"`
}

type UserRolesConfig struct {
	AWS map[string]*AWSRoleConfig `yaml:"aws"`
}

type UserConfig struct {
	PublicKey   string          `yaml:"public_key"`
	Roles       UserRolesConfig `yaml:"roles"`
	DefaultRole string          `yaml:"default_role"`
}

type ServerConfig struct {
	Port int `yaml:"port"`

	Users []*UserConfig `yaml:"users"`
}

type ClientConfig struct {
	ServerURL string   `yaml:"server_url"`
	KeyPath   string   `yaml:"key_path"`
	ShellPath string   `yaml:"shell_path"`
	Roles     []string `yaml:"roles"`
}

type Config struct {
	Client map[string]*ClientConfig `yaml:"client"`
	Server *ServerConfig            `yaml:"server"`
}
