package config_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/config"
)

const testConfig1Raw = `
target: https://cludo.superorbital.io/my-role
ssh_key_paths: ["~/.ssh/id_rsa_missing", "~/.ssh/id_ed25519"]
client:
  shell_path: "/usr/local/bin/bash"
server:
  port: 443
  github:
    api_endpoint: https://api.github.com/
  targets:
    prod:
      aws:
        arn: "aws:arn:iam:..."
        session_duration: "10m"
        access_key_id: "456DEF..."
        secret_access_key: "UVW789..."
    dev:
      aws:
        arn: "aws:arn:iam:..."
        session_duration: "8h"
        access_key_id: "123ABC..."
        secret_access_key: "ZXY098..."
  users:
  - public_key: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAII2YNmNYja2pH/D3hr8wwFqtRXIAdCYA25kgQESiWoDc test-key-ed25519-nopp@example.com DO NOT USE!!"
    name: "Aleema Bashir"
    github_id: "abashir"
    targets: ["prod", "dev"]
`

var testConfig1Duration1, _ = time.ParseDuration("10m")
var testConfig1Duration2, _ = time.ParseDuration("8h")
var testConfig1 = &config.Config{
	Target:      "https://cludo.superorbital.io/my-role",
	SSHKeyPaths: []string{"~/.ssh/id_rsa_missing", "~/.ssh/id_ed25519"},
	Client: &config.ClientConfig{
		ShellPath: "/usr/local/bin/bash",
	},
	Server: &config.ServerConfig{
		Github: &config.GithubConfig{
			APIEndpoint: "https://api.github.com/",
		},
		Port: 443,
		Users: []*config.UserConfig{
			{
				PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAII2YNmNYja2pH/D3hr8wwFqtRXIAdCYA25kgQESiWoDc test-key-ed25519-nopp@example.com DO NOT USE!!",
				Name:      "Aleema Bashir",
				GithubID:  "abashir",
				Targets:   []string{"prod", "dev"},
			},
		},
		Targets: map[string]*config.UserRolesConfig{
			"prod": {
				AWS: &config.AWSRoleConfig{
					AssumeRoleARN:   "aws:arn:iam:...",
					SessionDuration: testConfig1Duration1,
					AccessKeyID:     "456DEF...",
					SecretAccessKey: "UVW789...",
				},
			},
			"dev": {
				AWS: &config.AWSRoleConfig{
					AssumeRoleARN:   "aws:arn:iam:...",
					SessionDuration: testConfig1Duration2,
					AccessKeyID:     "123ABC...",
					SecretAccessKey: "ZXY098...",
				},
			},
		},
	},
}

func TestConfig(t *testing.T) {
	type test struct {
		name  string
		input string
		want  *config.Config
	}

	tests := []test{
		{
			name:  "Test config 1",
			input: testConfig1Raw,
			want:  testConfig1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile(".", "cludo*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(file.Name())

			_, err = file.WriteString(tc.input)
			if err != nil {
				t.Fatalf("Failed to populate temporary cludo config: %v", err)
			}

			viper.SetConfigFile(file.Name())
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					t.Fatal("ERROR: Failed to load configuration file: File not found")
				} else {
					t.Fatalf("ERROR: Failed to load configuration file: %v", err)
				}
			}

			got := &config.Config{}
			viper.Unmarshal(got)

			assert.EqualValuesf(t, tc.want, got, "expected: %#v, got: %#v", tc.want, got)
		})
	}
}
