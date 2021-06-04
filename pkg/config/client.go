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
)

type ClientConfig struct {
	Interactive bool   `mapstructure:"interactive"`
	KeyPath     string `mapstructure:"key_path"`
	Passphrase  string `mapstructure:"passphrase"`
	ShellPath   string `mapstructure:"shell_path"`
}

func (cc *ClientConfig) NewDefaultSigner() (*auth.Signer, error) {
	keyPath, err := homedir.Expand(cc.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand ~ (homedir) path characters in '%s': %v", cc.KeyPath, err)
	}
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

	var passphrase []byte = nil
	if cc.Passphrase != "" {
		passphrase = []byte(cc.Passphrase)
	}
	key, err := auth.DecodePrivateKey(encodedKey, passphrase, cc.Interactive)
	if err != nil {
		return nil, fmt.Errorf("Failed decoding private key %v: %v", keyPath, err)
	}

	return auth.NewDefaultSigner(key), nil
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
