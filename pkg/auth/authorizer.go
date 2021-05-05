package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Authorizer struct {
	users map[string]*rsa.PublicKey
}

func NewAuthorizer(users map[string]*rsa.PublicKey) *Authorizer {
	return &Authorizer{
		users: users,
	}
}

func (authz *Authorizer) CheckAuthHeader(message interface{}, header string) (string, bool, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return "", false, fmt.Errorf("Failed to JSON-serialize message for signature verification: %v", err)
	}
	// TODO: Make hashing/signing method pluggable.
	hashed := sha512.Sum512(b)
	hashedSlice := hashed[:]

	signature, err := base64.StdEncoding.DecodeString(header)
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
