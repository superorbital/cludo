package auth_test

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/auth"
	"golang.org/x/crypto/ssh"
)

func GenerateSSHAuthorizedKey(t *testing.T) (*rsa.PublicKey, []byte) {
	_, pubKey := GenerateRSAKeyPair(t)

	pub, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		t.Fatalf("Failed to convert pubkey to authorized key format: %v, %#v", err, pubKey)
	}

	return pubKey, ssh.MarshalAuthorizedKey(pub)
}

func GenerateSSHPrivateKey(t *testing.T) (*rsa.PrivateKey, []byte) {
	key, _ := GenerateRSAKeyPair(t)

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		t.Fatalf("Failed to convert private key to pem format: %v, %#v", err, key)
	}

	return key, private.Bytes()
}

func TestDecodeAuthorizedKey(t *testing.T) {
	type test struct {
		encoded []byte
		want    *rsa.PublicKey
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
		actual, actualErr := auth.DecodeAuthorizedKey(tc.encoded)

		assert.EqualValues(t, tc.want, actual)
		assert.EqualValues(t, tc.wantErr, actualErr)
	}
}

func TestEncodeAuthorizedKey(t *testing.T) {
	type test struct {
		pub     *rsa.PublicKey
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
		encoded []byte
		want    *rsa.PrivateKey
		wantErr error
	}

	key1, encoded1 := GenerateSSHPrivateKey(t)
	key2, encoded2 := GenerateSSHPrivateKey(t)

	tests := []test{
		{
			encoded: encoded1,
			want:    key1,
		},
		{
			encoded: encoded2,
			want:    key2,
		},
	}

	for _, tc := range tests {
		actual, actualErr := auth.DecodePrivateKey(tc.encoded, nil)

		assert.EqualValues(t, tc.want, actual)
		assert.EqualValues(t, tc.wantErr, actualErr)
	}
}

// TODO: Add test for private keys with passphrases
