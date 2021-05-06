package auth_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/auth"
)

func GenerateRSAKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("Failed to generate rsa keypair: %v", err)
	}
	return key, &key.PublicKey
}

func GenerateSHA512Signature(t *testing.T, key *rsa.PrivateKey, message string) string {
	b, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to serialize message: %v, %#v", err, message)
	}
	hashed := sha512.Sum512(b)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, hashed[:])
	if err != nil {
		t.Fatalf("Failed to generate signature: %v, %#v, %#v", err, key, hashed)
	}
	encoded := base64.StdEncoding.EncodeToString(signature)
	return fmt.Sprintf("%s:%s", message, encoded)
}

func TestAuthorizer(t *testing.T) {
	type test struct {
		message    string
		privateKey *rsa.PrivateKey
		allowed    map[string]*rsa.PublicKey
		want       string
		wantOK     bool
		wantErr    error
	}

	testKey1, testPub1 := GenerateRSAKeyPair(t)
	_, testPub2 := GenerateRSAKeyPair(t)
	testKey3, _ := GenerateRSAKeyPair(t)

	tests := []test{
		{
			message:    "test-message-1",
			privateKey: testKey1,
			allowed: map[string]*rsa.PublicKey{
				"test-id-1": testPub1,
			},
			want:   "test-id-1",
			wantOK: true,
		},
		{
			message:    "test-message-1",
			privateKey: testKey1,
			allowed: map[string]*rsa.PublicKey{
				"test-id-1": testPub1,
				"test-id-2": testPub2,
			},
			want:   "test-id-1",
			wantOK: true,
		},
		{
			message:    "test-message-1",
			privateKey: testKey3,
			allowed: map[string]*rsa.PublicKey{
				"test-id-1": testPub1,
				"test-id-2": testPub2,
			},
			want:   "",
			wantOK: false,
		},
		{
			message:    "test-message-1",
			privateKey: testKey3,
			allowed:    map[string]*rsa.PublicKey{},
			want:       "",
			wantOK:     false,
		},
	}

	for _, tc := range tests {
		authz := auth.NewAuthorizer(tc.allowed)
		actual, actualOK, actualErr := authz.CheckAuthHeader(GenerateSHA512Signature(t, tc.privateKey, tc.message))

		assert.EqualValues(t, tc.want, actual)
		assert.EqualValues(t, tc.wantOK, actualOK)
		assert.EqualValues(t, tc.wantErr, actualErr)
	}
}
