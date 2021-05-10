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
client:
  default:
    server_url: "https://www.example.com/"
    key_path: "~/.ssh/id_rsa"
    passphrase: ""
    shell_path: "/usr/local/bin/bash"
    roles: ["default"]
server:
  port: 443
  users:
    - public_key: "ssh-rsa aisudpoifueuyrlkjhflkyhaosiduyflakjsdhflkjashdf7898798765489..."
      roles:
        aws:
          so_org:
            arn: "aws:arn:iam:..."
            session_duration: "10m"
            access_key_id: "456DEF..."
            secret_access_key: "UVW789..."
          so_dev:
            arn: "aws:arn:iam:..."
            session_duration: "8h"
            access_key_id: "123ABC..."
            secret_access_key: "ZXY098..."
      default_role: "aws_so_org"
`

var testConfig1Duration1, _ = time.ParseDuration("10m")
var testConfig1Duration2, _ = time.ParseDuration("8h")
var testConfig1 = &config.Config{
	Client: map[string]*config.ClientConfig{
		config.DefaultClientConfig: {
			ServerURL:  "https://www.example.com/",
			KeyPath:    "~/.ssh/id_rsa",
			Passphrase: "",
			ShellPath:  "/usr/local/bin/bash",
			Roles:      []string{"default"},
		},
	},
	Server: &config.ServerConfig{
		Port: 443,
		Users: []*config.UserConfig{
			{
				PublicKey: "ssh-rsa aisudpoifueuyrlkjhflkyhaosiduyflakjsdhflkjashdf7898798765489...",
				Roles: &config.UserRolesConfig{
					AWS: map[string]*config.AWSRoleConfig{
						"so_org": {
							AssumeRoleARN:   "aws:arn:iam:...",
							SessionDuration: testConfig1Duration1,
							AccessKeyID:     "456DEF...",
							SecretAccessKey: "UVW789...",
						},
						"so_dev": {
							AssumeRoleARN:   "aws:arn:iam:...",
							SessionDuration: testConfig1Duration2,
							AccessKeyID:     "123ABC...",
							SecretAccessKey: "ZXY098...",
						},
					},
				},
				DefaultRole: "aws_so_org",
			},
		},
	},
}

func TestConfig(t *testing.T) {
	type test struct {
		input string
		want  *config.Config
	}

	tests := []test{
		{
			input: testConfig1Raw,
			want:  testConfig1,
		},
	}

	for _, tc := range tests {
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
	}
}
