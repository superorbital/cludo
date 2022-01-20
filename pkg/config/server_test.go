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
			PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC8zJJNfH4szhN+FU+5I8mtmbPU16F1ObcxY7rcwX9t5UxTMLB5PtTuw1LWNSp9b9XMHE4/Y/So9+P+CwgtjdyrfYeQ5aX+YeikK1+BC8Az0erY2JHfg8dsLJ8JGRw0SN7eYfJ/Kss4gTBF0NNFMTiity91i7A6yF+LcidjmYly4Qa0HWXxFcYpZ+u0Uj9BynkmDtJyfKWEBTqe/VlXV5N/tLOXvyotkPUlfSKv+d+6YOBVIMctlC0e7zIPxgG0UWr0ntzhzO9kYaOxmwelxjtE32rm+tt2RAs5JcxmgtppOn0SesvhgpF/iDt9TSZmqM0zc5FZjeY/ilQ5q7eMW0ZjP9kWeRnamJx0Cx5gNmPpUYLLdKOrEQNXT8FwkyfmRtQcgOkiPVFPXMnhHYO6DVAH0L5lGHL6jFHxX0SjfOEC63y7ehz7BFPVCxFpA12+HP93RzFV3d/ohaKobZkYa0qkMk+Nn2DAmV2msjGvwSS6QlFoUtbkcArqnqciTnNY3IF9IV+FRg2omzzoPA1rOLLthNsW84ycD58YzQrQ+8pkHpmDfspeSTr2Jges5E1Z6koBdOaeC6p/Ud6EnDG9Plo7I1yeBaYKh1zZEVO7L6GvMB17xqSk5sXIj0AupavGqEiju10SnnnDfZ1aroMPWKT+4aaE5WEarzYeuqiouDxDoQ== test-key-nopp@example.com DO NOT USE!!",
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
