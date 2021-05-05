package auth_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/pkg/auth"
)

func VerifyHeader(t *testing.T, message []byte, header string, publicKey *rsa.PublicKey) {
	hashed := sha512.Sum512(message)

	signature, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		t.Fatalf("Failed to decode header: %v, %#v", err, header)
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashed[:], signature)
	assert.NoErrorf(t, err, "Failed to verify generated header: %v, %#v", err, header)
}

func TestSigner(t *testing.T) {
	type test struct {
		message    []byte
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		want       bool
		wantErr    error
	}

	testKey1, testPub1 := GenerateRSAKeyPair(t)
	testKey2, _ := GenerateRSAKeyPair(t)

	tests := []test{
		{
			message:    []byte("test-message-1"),
			privateKey: testKey1,
			publicKey:  testPub1,
		},
		{
			message:    []byte("test-message-2"),
			privateKey: testKey2,
			publicKey:  testPub1,

			want: false,
		},
	}

	for _, tc := range tests {
		signer := auth.NewDefaultSigner(tc.privateKey)
		actual, actualErr := signer.GenerateAuthHeader(tc.message)

		assert.EqualValues(t, tc.wantErr, actualErr)
		if tc.want {
			VerifyHeader(t, tc.message, actual, tc.publicKey)
		}
	}
}
