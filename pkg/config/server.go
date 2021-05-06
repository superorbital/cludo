package config

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/superorbital/cludo/pkg/auth"
)

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

func (uc *UserConfig) ID() string {
	return base64.StdEncoding.EncodeToString([]byte(uc.PublicKey))
}

type ServerConfig struct {
	Port int `yaml:"port"`

	Users []*UserConfig `yaml:"users"`
}

func (sc *ServerConfig) NewAuthorizer() (*auth.Authorizer, error) {
	users := map[string]*rsa.PublicKey{}
	for _, user := range sc.Users {
		pub, err := auth.DecodeAuthorizedKey([]byte(user.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("Failed to decode user public key: %v, %#v", err, user.PublicKey)
		}
		users[user.ID()] = pub
	}
	return auth.NewAuthorizer(users), nil
}
