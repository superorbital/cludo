package config

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/superorbital/cludo/pkg/auth"
	"github.com/superorbital/cludo/pkg/aws"
)

type AWSRoleConfig struct {
	SessionDuration time.Duration `mapstructure:"session_duration"`
	AccessKeyID     string        `mapstructure:"access_key_id"`
	SecretAccessKey string        `mapstructure:"secret_access_key"`
	AssumeRoleARN   string        `mapstructure:"arn"`
	Region          string        `mapstructure:"region"`
}

func (arc *AWSRoleConfig) NewPlugin() (*aws.AWSPlugin, error) {
	return aws.NewAWSPlugin(arc.AccessKeyID, arc.SecretAccessKey, arc.Region, arc.SessionDuration, arc.AssumeRoleARN)
}

type UserRolesConfig struct {
	AWS *AWSRoleConfig `mapstructure:"aws"`
}

type UserConfig struct {
	PublicKey string   `mapstructure:"public_key"`
	Name      string   `mapstructure:"name"`
	Targets   []string `mapstructure:"targets"`
}

func (uc *UserConfig) ID() string {
	return base64.StdEncoding.EncodeToString([]byte(uc.PublicKey))
}

type ServerConfig struct {
	Port    int                         `yaml:"port"`
	Targets map[string]*UserRolesConfig `mapstructure:"targets"`

	Users []*UserConfig `mapstructure:"users"`
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

func (sc *ServerConfig) GetUser(id string) (*UserConfig, bool) {
	for _, user := range sc.Users {
		if user.ID() == id {
			return user, true
		}
	}
	return nil, false
}
