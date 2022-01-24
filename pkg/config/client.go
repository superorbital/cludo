package config

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/mitchellh/go-homedir"
	"github.com/superorbital/cludo/client"
	"github.com/superorbital/cludo/pkg/auth"
	"golang.org/x/crypto/ssh"
)

type ClientConfig struct {
	Interactive bool   `mapstructure:"interactive"`
	ShellPath   string `mapstructure:"shell_path"`
}

// NewDefaultClientSigner attempts to read and decode the provided private key
// and then generate a signer that can be used to sign requests to the server.
// It returns a *auth.Signer and any errors that were encountered.
func (cc *ClientConfig) NewDefaultClientSigner(pkey string) (*auth.Signer, error) {
	keyPath, err := homedir.Expand(pkey)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand ~ (homedir) path characters in '%s': %v", pkey, err)
	}
	if keyPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Failed to get user home directory: %v", err)
		}

		// Set keypath to ~/.ssh/id_rsa by default if nothing is set.
		keyPath = path.Join(homeDir, ".ssh", "id_rsa")
	}

	rawKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read user key %v: %v", keyPath, err)
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(rawKey))
	if err != nil {
		// This is not a public key, so it might be a private key

		key, err := auth.DecodePrivateKey(keyPath, rawKey, cc.Interactive)
		if err != nil {
			return nil, fmt.Errorf("Failed decoding key %v: %v", keyPath, err)
		}
		return auth.NewDefaultSigner(key, nil), nil
	}
	// This is a public key, so we pass it through
	return auth.NewDefaultSigner(nil, publicKey), nil
}

func NewClient(target string, debug bool) (*client.Cludod, error) {
	serverURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse client target URL '%s': %v", target, err)
	}

	r := httptransport.New(
		serverURL.Host,
		path.Join(client.DefaultBasePath),
		[]string{serverURL.Scheme},
	)
	r.SetDebug(debug)

	// Set custom producer and consumer to use the default ones
	r.Consumers["application/json"] = runtime.JSONConsumer()
	r.Producers["application/json"] = runtime.JSONProducer()

	return client.New(r, strfmt.Default), nil
}
