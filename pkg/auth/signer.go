package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const CludoAuthHeader = "X-CLUDO-KEY"

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
	// TODO: Make hashing/signing method pluggable.
	hashed := sha512.Sum512([]byte(message))
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
	// Use base64 encoding to prevent the use of non-displayable characters.
	message := base64.StdEncoding.EncodeToString(b)
	return signer.GenerateAuthHeader(message)
}

// CludoAuth provides an API key auth info writer
func (signer *Signer) CludoAuth() runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		value, err := signer.GenerateRandomAuthHeader()
		if err != nil {
			return fmt.Errorf("Failed to generate cludo auth header: %v", err)
		}
		return r.SetHeaderParam(CludoAuthHeader, value)
	})
}
