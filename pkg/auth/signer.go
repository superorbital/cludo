package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const CludoAuthHeader = "X-CLUDO-KEY"

type Signer struct {
	rng        io.Reader
	privateKey *rsa.PrivateKey
	publicKey  ssh.PublicKey
}

// NewDefaultSigner runs the NewSigner() function
// It return the response from NewSigner()
func NewDefaultSigner(privkey *rsa.PrivateKey, pubkey ssh.PublicKey) *Signer {
	return NewSigner(privkey, pubkey, rand.Reader)
}

// NewSigner creates a new Signer that can be used to sign requests
// It returns a pointer to the signer.
func NewSigner(privkey *rsa.PrivateKey, pubkey ssh.PublicKey, rng io.Reader) *Signer {
	return &Signer{
		rng:        rng,
		privateKey: privkey,
		publicKey:  pubkey,
	}
}

func (signer *Signer) GenerateAuthHeader(message string) (string, error) {
	// TODO: Make hashing/signing method pluggable.
	var encoded string
	if signer.publicKey != nil {
		/*
			If we have a publicKey then we want to try and
			sign the message via a SSH Auth Socket

			Note: This uses SHA1 versus SHA512
			as the SSH Auth Socket spec only supports SHA1
		*/
		socket := os.Getenv("SSH_AUTH_SOCK")
		if socket != "" {
			conn, err := net.Dial("unix", socket)
			if err != nil {
				return "", fmt.Errorf("Failed to open SSH_AUTH_SOCK: %v", err)
			}
			client := agent.NewClient(conn)
			sig, err := client.Sign(signer.publicKey, []byte(message))
			if err != nil {
				return "", fmt.Errorf("failed to sign message via SSH Auth Socket: %v", err)
			}
			encoded = "sha1|" + base64.StdEncoding.EncodeToString(sig.Blob)

		} else {
			return "", fmt.Errorf("Received SSH public key, but SSH_AUTH_SOCK is not set in the environment")
		}
	} else {
		// If we don't have a publicKey then we want to sign the message via a decoded private key
		hashed := sha512.Sum512([]byte(message))

		signature, err := rsa.SignPKCS1v15(signer.rng, signer.privateKey, crypto.SHA512, hashed[:])
		encoded = "sha512|" + base64.StdEncoding.EncodeToString(signature)
		if err != nil {
			return "", fmt.Errorf("failed to sign message via SSH private key: %v", err)
		}
	}

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
