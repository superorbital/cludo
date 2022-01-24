package auth_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/auth"
)

func VerifyHeader(t *testing.T, header string, publicKey *rsa.PublicKey) {
	// API tokens are of the form: <random-number>:<signature-of-random-number>
	tokens := strings.SplitN(header, ":", 2)
	if len(tokens) != 2 {
		t.Fatalf("Found malformed auth header: %v", header)
	}

	hashed := sha512.Sum512([]byte(tokens[0]))

	signature, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		t.Fatalf("Failed to decode header: %v, %#v", err, tokens[1])
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashed[:], signature)
	assert.NoErrorf(t, err, "Failed to verify generated header: %v, %#v", err, header)
}

func TestSigner(t *testing.T) {
	type test struct {
		name       string
		message    string
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		want       bool
		wantErr    error
	}

	testKey1, testPub1 := GenerateRSAKeyPair(t)
	testKey2, _ := GenerateRSAKeyPair(t)

	tests := []test{
		{
			name:       "Test matching keys",
			message:    "test-message-1",
			privateKey: testKey1,
			publicKey:  testPub1,
		},
		{
			name:       "Test non-matching keys",
			message:    "test-message-2",
			privateKey: testKey2,
			publicKey:  testPub1,

			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			signer := auth.NewDefaultSigner(tc.privateKey, nil)
			actual, actualErr := signer.GenerateAuthHeader(tc.message)

			assert.EqualValues(t, tc.wantErr, actualErr)
			if tc.want {
				VerifyHeader(t, actual, tc.publicKey)
			}
		})
	}
}

// TODO: Create tests for SSH Public Key/Agent workflow
