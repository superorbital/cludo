package config

import (
	"crypto/rsa"
	"encoding/base64"
	"log"
	"time"

	"github.com/superorbital/cludo/pkg/auth"
	"github.com/superorbital/cludo/pkg/aws"
	"github.com/superorbital/cludo/pkg/providers"
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
	GithubID  string   `mapstructure:"github_id"`
	Targets   []string `mapstructure:"targets"`
}

func (uc *UserConfig) ID() string {
	return base64.StdEncoding.EncodeToString([]byte(uc.PublicKey))
}

type GithubConfig struct {
	APIEndpoint string `mapstructure:"api_endpoint"`
}

type ServerConfig struct {
	Port    int                         `yaml:"port"`
	Targets map[string]*UserRolesConfig `mapstructure:"targets"`

	Github *GithubConfig `mapstructure:"github"`
	Users  []*UserConfig `mapstructure:"users"`
}

func (sc *ServerConfig) NewConfigAuthorizer() (*auth.Authorizer, error) {
	users := map[string]*rsa.PublicKey{}
	for _, user := range sc.Users {
		pub, err := auth.DecodeAuthorizedKey([]byte(user.PublicKey))
		if err != nil {
			log.Printf("[WARN] Failed to decode user public key: %v, %#v\n", err, user.PublicKey)
		}
		users[user.ID()] = pub
	}
	return auth.NewAuthorizer(users), nil
}

func (sc *ServerConfig) NewGithubAuthorizer() (*auth.Authorizer, error) {
	users := map[string]*rsa.PublicKey{}

	for _, user := range sc.Users {
		if user.GithubID != "" {
			api_endpoint := ""
			if sc.Github != nil {
				api_endpoint = sc.Github.APIEndpoint
			}
			provider_keys, err := providers.CollectGithubPublicKeys(api_endpoint, user.GithubID)
			if err != nil {
				return nil, err
			}
			for _, pkey := range provider_keys {
				pub, err := auth.DecodeAuthorizedKey([]byte(pkey))
				if err != nil {
					log.Printf("[WARN] Failed to decode user public key: %v, %#v\n", err, pkey)
				} else {
					users[user.ID()] = pub
				}
			}
		}
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
