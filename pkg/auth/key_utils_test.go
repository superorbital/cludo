package auth_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/auth"
	"golang.org/x/crypto/ssh"
)

func GenerateSSHAuthorizedKey(t *testing.T) (*ssh.PublicKey, []byte) {
	_, pubKey := GenerateRSAKeyPair(t)

	pub, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		t.Fatalf("Failed to convert pubkey to authorized key format: %v, %#v", err, pubKey)
	}

	return &pub, ssh.MarshalAuthorizedKey(pub)
}

func GenerateSSHRSAPrivateKey(t *testing.T, passphrase string) (*rsa.PrivateKey, []byte) {
	key, _ := GenerateRSAKeyPair(t)

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	var err error
	// Encrypt the PEM
	if passphrase != "" {
		privateKeyPEM, err = x509.EncryptPEMBlock(rand.Reader, privateKeyPEM.Type, privateKeyPEM.Bytes, []byte(passphrase), x509.PEMCipherAES256)
		if err != nil {
			return nil, nil
		}
	}

	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		t.Fatalf("Failed to convert private key to pem format: %v, %#v", err, key)
	}

	return key, private.Bytes()
}

func TestParseAuthorizedKey(t *testing.T) {
	type test struct {
		encoded []byte
		want    *ssh.PublicKey
		wantErr error
	}

	testPub1, encoded1 := GenerateSSHAuthorizedKey(t)
	testPub2, encoded2 := GenerateSSHAuthorizedKey(t)

	tests := []test{
		{
			encoded: encoded1,
			want:    testPub1,
		},
		{
			encoded: encoded2,
			want:    testPub2,
		},
	}

	for _, tc := range tests {
		actual, actualErr := auth.ParseAuthorizedKey(tc.encoded)

		assert.EqualValues(t, tc.want, actual)
		assert.EqualValues(t, tc.wantErr, actualErr)
	}
}

func TestEncodeAuthorizedKey(t *testing.T) {
	type test struct {
		pub     *ssh.PublicKey
		want    []byte
		wantErr error
	}

	testPub1, encoded1 := GenerateSSHAuthorizedKey(t)
	testPub2, encoded2 := GenerateSSHAuthorizedKey(t)

	tests := []test{
		{
			pub:  testPub1,
			want: encoded1,
		},
		{
			pub:  testPub2,
			want: encoded2,
		},
	}

	for _, tc := range tests {
		actual, actualErr := auth.EncodeAuthorizedKey(tc.pub)

		assert.EqualValues(t, tc.want, actual)
		assert.EqualValues(t, tc.wantErr, actualErr)
	}
}

func TestDecodePrivateKey(t *testing.T) {
	type test struct {
		name    string
		encoded []byte
		want    *interface{}
		wantErr error
	}

	passphrase := ""
	_, encoded1 := GenerateSSHRSAPrivateKey(t, passphrase)
	_, encoded2 := GenerateSSHRSAPrivateKey(t, passphrase)

	tests := []test{
		{
			name:    "Test key 1",
			encoded: encoded1,
		},
		{
			name:    "Test key 2",
			encoded: encoded2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := auth.DecodePrivateKey(strings.Replace(strings.ToLower(tc.name), " ", "-", -1), tc.encoded, false)

			assert.NoErrorf(t, err, "Failed to decode private key: %v", err)
		})
	}
}

func TestDecodePrivateKeyWithPassPhrase(t *testing.T) {
	type test struct {
		name       string
		encoded    []byte
		passphrase string
		want       *interface{}
		wantErr    error
	}

	passphrase1 := "cludo123"
	_, encoded1 := GenerateSSHRSAPrivateKey(t, passphrase1)

	tests := []test{
		{
			name:       "Test passphrase key 1",
			encoded:    encoded1,
			passphrase: "badpassphrase",
			want:       nil,
			wantErr:    errors.New("Failed to parse private key: ssh: this private key is passphrase protected"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := auth.DecodePrivateKey(strings.Replace(strings.ToLower(tc.name), " ", "-", -1), tc.encoded, false)

			assert.EqualErrorf(t, err, tc.wantErr.Error(), "did not get expected error", "formatted")
		})
	}
}
