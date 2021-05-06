package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strings"
)

type Authorizer struct {
	users map[string]*rsa.PublicKey
}

func NewAuthorizer(users map[string]*rsa.PublicKey) *Authorizer {
	return &Authorizer{
		users: users,
	}
}

func (authz *Authorizer) CheckAuthHeader(header string) (string, bool, error) {
	// API tokens are of the form: <random-number>:<signature-of-random-number>
	tokens := strings.SplitN(header, ":", 2)
	if len(tokens) != 2 {
		return "", false, fmt.Errorf("Found malformed auth header: %v", header)
	}

	// TODO: Make hashing/signing method pluggable.
	hashed := sha512.Sum512([]byte(tokens[0]))
	hashedSlice := hashed[:]

	signature, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return "", false, fmt.Errorf("Failed to base64-decode auth header for signature verification: %v", err)
	}

	for id, publicKey := range authz.users {
		// TODO: Make hashing/signing method pluggable.
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashedSlice, signature)
		if err == nil {
			return id, true, nil
		}
	}

	return "", false, nil
}
