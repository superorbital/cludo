package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/superorbital/cludo/pkg/auth"
)

type ClientConfig struct {
	ServerURL  string   `mapstructure:"server_url"`
	KeyPath    string   `mapstructure:"key_path"`
	Passphrase string   `mapstructure:"passphrase"`
	ShellPath  string   `mapstructure:"shell_path"`
	Roles      []string `mapstructure:"roles"`
}

func (cc *ClientConfig) NewDefaultSigner() (*auth.Signer, error) {
	keyPath := cc.KeyPath
	if keyPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Failed to get user home directory: %v", err)
		}

		// Set keypath to ~/.ssh/id_rsa by default if nothing is set.
		keyPath = path.Join(homeDir, ".ssh", "id_rsa")
	}

	encodedKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read user private key %v: %v", keyPath, err)
	}

	// TODO: Handle passwords.
	key, err := auth.DecodePrivateKey(encodedKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed decoding private key %v: %v", keyPath, err)
	}

	return auth.NewDefaultSigner(key), nil
}
