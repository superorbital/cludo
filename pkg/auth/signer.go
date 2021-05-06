package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type Signer struct {
	rng        io.Reader
	privateKey *rsa.PrivateKey
}

func NewDefaultSigner(key *rsa.PrivateKey) *Signer {
	return NewSigner(key, rand.Reader)
}

func NewSigner(key *rsa.PrivateKey, rng io.Reader) *Signer {
	return &Signer{
		rng:        rng,
		privateKey: key,
	}
}

func (signer *Signer) GenerateAuthHeader(message string) (string, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("Failed to JSON-serialize message for signing: %v", err)
	}

	// TODO: Make hashing/signing method pluggable.
	hashed := sha512.Sum512(b)
	signature, err := rsa.SignPKCS1v15(signer.rng, signer.privateKey, crypto.SHA512, hashed[:])
	if err != nil {
		return "", fmt.Errorf("Failed to sign message: %v:", err)
	}

	encoded := base64.StdEncoding.EncodeToString(signature)
	return fmt.Sprintf("%s:%s", message, encoded), nil
}

func (signer *Signer) GenerateRandomAuthHeader() (string, error) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Failed to generate random message: %v", err)
	}
	message := string(b)
	return signer.GenerateAuthHeader(message)
}
