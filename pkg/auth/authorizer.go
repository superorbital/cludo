package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
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

// TODO: This is still subject to replay attacks unless we somehow prevent old messages from being re-used.
func (authz *Authorizer) CheckAuthHeader(header string) (string, bool, error) {
	// API tokens are of the form: <random-number>:<signature-of-random-number>
	tokens := strings.SplitN(header, ":", 2)
	if len(tokens) != 2 {
		return "", false, fmt.Errorf("Found malformed auth header: %v", header)
	}

	hashSignature := strings.SplitN(tokens[1], "|", 2)
	if len(hashSignature) != 2 {
		return "", false, fmt.Errorf("Found malformed auth header signature: %v", header)
	}
	signature, err := base64.StdEncoding.DecodeString(hashSignature[1])
	if err != nil {
		return "", false, fmt.Errorf("Failed to base64-decode auth header for signature verification: %v", err)
	}
	if hashSignature[0] == "sha1" {
		// TODO: Make hashing/signing method pluggable.
		hashed := sha1.Sum([]byte(tokens[0]))
		hashedSlice := hashed[:]
		for id, publicKey := range authz.users {
			err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedSlice, signature)
			if err == nil {
				return id, true, nil
			}
		}
	} else if hashSignature[0] == "sha512" {
		// TODO: Make hashing/signing method pluggable.
		hashed := sha512.Sum512([]byte(tokens[0]))
		hashedSlice := hashed[:]
		for id, publicKey := range authz.users {
			err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashedSlice, signature)
			if err == nil {
				return id, true, nil
			}
		}
	} else {
		return "", false, fmt.Errorf("Unknown signature hashing prefix: %v", hashSignature[0])
	}

	return "", false, nil
}
