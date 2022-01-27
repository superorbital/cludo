package auth_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/auth"
	"golang.org/x/crypto/ssh"
)

func GenerateRSAKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("Failed to generate rsa keypair: %v", err)
	}

	return key, &key.PublicKey
}

func PrivateKeyRSAToInterface(t *testing.T, key *rsa.PrivateKey) (*interface{}, error) {
	privateKeyDer := x509.MarshalPKCS1PrivateKey(key)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privateKeyPem := pem.EncodeToMemory(&privateKeyBlock)

	keyInterface, err := ssh.ParseRawPrivateKey(privateKeyPem)
	if err != nil {
		return nil, fmt.Errorf("Failed decoding key: %v", err)
	}
	return &keyInterface, nil
}

func GenerateSHA512Signature(t *testing.T, key *rsa.PrivateKey, message string) string {
	hashed := sha512.Sum512([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, hashed[:])
	if err != nil {
		t.Fatalf("Failed to generate signature: %v, %#v, %#v", err, key, hashed)
	}
	encoded := "direct|" + base64.StdEncoding.EncodeToString(signature)
	return fmt.Sprintf("%s:%s", message, encoded)
}

func TestAuthorizer(t *testing.T) {
	type test struct {
		name       string
		message    string
		privateKey *rsa.PrivateKey
		allowed    map[string]*ssh.PublicKey
		want       string
		wantOK     bool
		wantErr    error
	}

	testKey1, testPub1 := GenerateRSAKeyPair(t)
	_, testPub2 := GenerateRSAKeyPair(t)
	testKey3, _ := GenerateRSAKeyPair(t)
	tpk1, err := ssh.NewPublicKey(testPub1)
	if err != nil {
		t.Fatalf("Failed to generate public key (tpk1): %v", err)
	}
	tpk2, err := ssh.NewPublicKey(testPub2)
	if err != nil {
		t.Fatalf("Failed to generate public key (tpk2): %v", err)
	}

	tests := []test{
		{
			name:       "Test sole matching key",
			message:    "test-message-1",
			privateKey: testKey1,
			allowed: map[string]*ssh.PublicKey{
				"test-id-1": &tpk1,
			},
			want:   "test-id-1",
			wantOK: true,
		},
		{
			name:       "Test matching key",
			message:    "test-message-1",
			privateKey: testKey1,
			allowed: map[string]*ssh.PublicKey{
				"test-id-1": &tpk1,
				"test-id-2": &tpk2,
			},
			want:   "test-id-1",
			wantOK: true,
		},
		{
			name:       "Test non-matching key",
			message:    "test-message-1",
			privateKey: testKey3,
			allowed: map[string]*ssh.PublicKey{
				"test-id-1": &tpk1,
				"test-id-2": &tpk2,
			},
			want:   "",
			wantOK: false,
		},
		{
			name:       "Test empty authorizer",
			message:    "test-message-1",
			privateKey: testKey3,
			allowed:    map[string]*ssh.PublicKey{},
			want:       "",
			wantOK:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authz := auth.NewAuthorizer(tc.allowed)
			actual, actualOK, actualErr := authz.CheckAuthHeader(GenerateSHA512Signature(t, tc.privateKey, tc.message))

			assert.EqualValues(t, tc.want, actual)
			assert.EqualValues(t, tc.wantOK, actualOK)
			assert.EqualValues(t, tc.wantErr, actualErr)
		})
	}
}
