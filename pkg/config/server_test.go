package config

import (
	"testing"
	"time"
)

var testConfig1Duration, _ = time.ParseDuration("1h")
var testConfig1 = &ServerConfig{
	Github: &GithubConfig{
		APIEndpoint: "https://api.github.com/",
	},
	Port: 443,
	Users: []*UserConfig{
		{
			PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAII2YNmNYja2pH/D3hr8wwFqtRXIAdCYA25kgQESiWoDc test-key-ed25519-nopp@example.com DO NOT USE!!",
			Name:      "Aleema Bashir",
			GithubID:  "abashir",
			Targets:   []string{"prod", "dev"},
		},
	},
	Targets: map[string]*UserRolesConfig{
		"prod": {
			AWS: &AWSRoleConfig{
				AssumeRoleARN:   "aws:arn:iam:...",
				SessionDuration: testConfig1Duration,
				AccessKeyID:     "456DEF...",
				SecretAccessKey: "UVW789...",
			},
		},
		"dev": {
			AWS: &AWSRoleConfig{
				AssumeRoleARN:   "aws:arn:iam:...",
				SessionDuration: testConfig1Duration,
				AccessKeyID:     "123ABC...",
				SecretAccessKey: "ZXY098...",
			},
		},
	},
}
var testAuthorizer1 = ""

func TestConfigAuthorizer(t *testing.T) {
	type test struct {
		name  string
		input *ServerConfig
	}

	tests := []test{
		{
			name:  "Test ConfigAuthorizer 1",
			input: testConfig1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.NewConfigAuthorizer()
			if err != nil {
				t.Fatalf("Failed to populate authorizer from config: %v", err)
			}
		})
	}
}

func TestGithubAuthorizer(t *testing.T) {
	type test struct {
		name  string
		input *ServerConfig
	}

	tests := []test{
		{
			name:  "Test GithubAuthorizer 1",
			input: testConfig1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.NewGithubAuthorizer()
			if err != nil {
				t.Fatalf("Failed to populate authorizer from config: %v", err)
			}
		})
	}
}
