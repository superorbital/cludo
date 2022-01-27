package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
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
	privateKey *interface{}
	publicKey  *ssh.PublicKey
}

// NewDefaultSigner runs the NewSigner() function
// It return the response from NewSigner()
func NewDefaultSigner(privkey *interface{}, pubkey *ssh.PublicKey) *Signer {
	return NewSigner(privkey, pubkey, rand.Reader)
}

// NewSigner creates a new Signer that can be used to sign requests
// It returns a pointer to the signer.
func NewSigner(privkey *interface{}, pubkey *ssh.PublicKey, rng io.Reader) *Signer {
	return &Signer{
		rng:        rng,
		privateKey: privkey,
		publicKey:  pubkey,
	}
}

func (signer *Signer) GenerateAuthHeader(message string) (string, error) {
	// TODO: Make hashing/signing method pluggable.
	var encoded string
	var err error
	if signer.publicKey != nil {
		/*
			If we have a publicKey then we want to try and
			sign the message via a SSH Auth Socket

			Note:
				For RSA the SSH agent uses SHA1 for hashing.
				For ECDSA the SSH agent uses SHA256 for hashing and a non-ASN1 response.
		*/
		socket := os.Getenv("SSH_AUTH_SOCK")
		if socket != "" {
			encoded, err = signWithAgent(socket, message, *signer.publicKey)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("Received SSH public key, but SSH_AUTH_SOCK is not set in the environment")
		}
	} else {
		// If we don't have a publicKey then we want to sign the message via a decoded private key
		var keyType string
		var signature []byte
		var err error
		privKey := *signer.privateKey
		keyType, err = detectPrivateKeyType(signer.privateKey)
		if err != nil {
			return "", fmt.Errorf("failed to sign message: %v", err)
		}
		if keyType == "rsa" {
			hashed := sha512.Sum512([]byte(message))
			signature, err = rsa.SignPKCS1v15(signer.rng, privKey.(*rsa.PrivateKey), crypto.SHA512, hashed[:])
			if err != nil {
				return "", fmt.Errorf("failed to sign message via SSH RSA private key: %v", err)
			}
			encoded = "direct|" + base64.StdEncoding.EncodeToString(signature)
		} else if keyType == "ed25519" {
			ed25519Key, ok := privKey.(*ed25519.PrivateKey)
			if !ok {
				return "", fmt.Errorf("failed to convert ed25519 private key: %v", err)
			}
			signature = ed25519.Sign(*ed25519Key, []byte(message))
			encoded = "direct|" + base64.StdEncoding.EncodeToString(signature)
		} else if keyType == "ecdsa" {
			hashed := []byte(message)
			signature, err = ecdsa.SignASN1(signer.rng, privKey.(*ecdsa.PrivateKey), hashed[:])
			if err != nil {
				return "", fmt.Errorf("failed to sign message via SSH ECDSA private key: %v", err)
			}
			encoded = "direct|" + base64.StdEncoding.EncodeToString(signature)
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

func signWithAgent(socket string, message string, pkey ssh.PublicKey) (string, error) {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return "", fmt.Errorf("Failed to open SSH_AUTH_SOCK: %v", err)
	}
	client := agent.NewClient(conn)
	sig, err := client.Sign(pkey, []byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to sign message via SSH Auth Socket: %v", err)
	}
	return "agent|" + base64.StdEncoding.EncodeToString(sig.Blob), nil
}

func detectPrivateKeyType(privKey *interface{}) (string, error) {
	keyType := ""

	if isRSAPrivateKey(privKey) == true {
		keyType = "rsa"
	} else if isED25519PrivateKey(privKey) == true {
		keyType = "ed25519"
	} else if isECDSAPrivateKey(privKey) == true {
		keyType = "ecdsa"
	} else {
		// TODO: Fix error passing upstream.
		// We are clobbering many of these messages
		return "", fmt.Errorf("Unrecognized private key type")
	}
	return keyType, nil
}

func isRSAPrivateKey(privKey *interface{}) bool {
	key := *privKey
	_, ok := key.(*rsa.PrivateKey)
	if ok {
		return true
	}
	return false
}

func isED25519PrivateKey(privKey *interface{}) bool {
	key := *privKey
	_, ok := key.(*ed25519.PrivateKey)
	if ok {
		return true
	}
	return false
}

func isECDSAPrivateKey(privKey *interface{}) bool {
	key := *privKey
	_, ok := key.(*ecdsa.PrivateKey)
	if ok {
		return true
	}
	return false
}
