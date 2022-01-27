package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strings"

	"golang.org/x/crypto/ssh"
)

/*
  TODO: There is some good logic and techniques in here that we might be able
	to use to make this all more flexible and "correct".

	Determining the hash type, etc.

	https://github.com/golang/crypto/blob/master/ssh/keys.go
*/

type Authorizer struct {
	users map[string]*ssh.PublicKey
}

// NewAuthorizer creates a new Authorizer struct.
// It returns a pointer to the Authorizer.
func NewAuthorizer(users map[string]*ssh.PublicKey) *Authorizer {
	return &Authorizer{
		users: users,
	}
}

// CheckAuthHeader checks the given signature against the avaliable public keys.
// It returns the user ID if the signature is valid, in addition to a
// verification boolean and error.
func (authz *Authorizer) CheckAuthHeader(header string) (string, bool, error) {
	// TODO: This is still subject to replay attacks unless we somehow prevent
	// old messages from being re-used.
	//
	// API tokens are of the form: <random-number>:<signature-of-random-number>
	tokens := strings.SplitN(header, ":", 2)
	if len(tokens) != 2 {
		return "", false, fmt.Errorf("Found malformed auth header: %v", header)
	}

	hashSignature := strings.SplitN(tokens[1], "|", 2)
	if len(hashSignature) != 2 {
		return "", false, fmt.Errorf("Found malformed auth header signature: %v", header)
	}
	hashSource := hashSignature[0]
	signature, err := base64.StdEncoding.DecodeString(hashSignature[1])
	if err != nil {
		return "", false, fmt.Errorf("Failed to base64-decode auth header for signature verification: %v", err)
	}
	if hashSource == "agent" {
		// TODO: Make hashing/signing method pluggable.
		for id, publicKey := range authz.users {
			pubKey := *publicKey
			keyType := pubKey.Type()

			switch keyType {
			case "ssh-rsa":
				hash := "sha1"
				err := VerifyRSAHeader(pubKey, tokens, signature, hash)
				if err == nil {
					return id, true, nil
				}
			case "ssh-ed25519":
				err := VerifyED25519Header(pubKey, tokens, signature)
				if err == nil {
					return id, true, nil
				}
			case "ecdsa-sha2-nistp256":
				sigType := "default"
				verified, err := VerifyECDSAHeader(pubKey, tokens, signature, sigType)
				if verified == true && err == nil {
					return id, true, nil
				}
			default:
				log.Printf("[WARN] %s is a key type that is not supported", keyType)
			}
		}
	} else if hashSource == "direct" {
		// TODO: Make hashing/signing method pluggable.
		for id, publicKey := range authz.users {
			pubKey := *publicKey
			keyType := pubKey.Type()

			switch keyType {
			case "ssh-rsa":
				hash := "sha512"
				VerifyRSAHeader(pubKey, tokens, signature, hash)
				if err == nil {
					return id, true, nil
				}
			case "ssh-ed25519":
				VerifyED25519Header(pubKey, tokens, signature)
				if err == nil {
					return id, true, nil
				}
			case "ecdsa-sha2-nistp256":
				sigType := "asn1"
				verified, err := VerifyECDSAHeader(pubKey, tokens, signature, sigType)
				if verified == true && err == nil {
					return id, true, nil
				}
			default:
				log.Printf("[WARN] %s is a key type that is not supported", keyType)
			}
		}
	}
	return "", false, nil
}

// VerifyRSAHeader verifies a signature using the RSA algorithm.
// It returns nil if the signature is valid, or an error if it is invalid.
func VerifyRSAHeader(pubKey interface{}, tokens []string, signature []byte, hash string) error {
	parsedCryptoKey := pubKey.(ssh.CryptoPublicKey)
	pubCrypto := parsedCryptoKey.CryptoPublicKey()
	pkey, ok := pubCrypto.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("expected an RSA public key, got %T", pkey)
	}
	switch hash {
	case "sha1":
		hashed := sha1.Sum([]byte(tokens[0]))
		hashedSlice := hashed[:]
		err := rsa.VerifyPKCS1v15(pkey, crypto.SHA1, hashedSlice, signature)
		if err != nil {
			return err
		}
	case "sha512":
		hashed := sha512.Sum512([]byte(tokens[0]))
		hashedSlice := hashed[:]
		err := rsa.VerifyPKCS1v15(pkey, crypto.SHA1, hashedSlice, signature)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported hash type for RSA signature: %s", hash)
	}
	return nil
}

// VerifyED25519Header verifies a signature using the Ed25519 algorithm.
// It returns nil if the signature is valid, or an error if it is invalid.
func VerifyED25519Header(pubKey interface{}, tokens []string, signature []byte) error {
	parsedCryptoKey := pubKey.(ssh.CryptoPublicKey)
	pubCrypto := parsedCryptoKey.CryptoPublicKey()
	pkey, ok := pubCrypto.(ed25519.PublicKey)
	if !ok {
		return fmt.Errorf("expected an ed25519 public key, got %T", pkey)
	}
	hashed := []byte(tokens[0])
	hashedSlice := hashed[:]
	verified := ed25519.Verify(pkey, hashedSlice, signature)
	if verified != true {
		return fmt.Errorf("unable to verify ed25519 signature")
	}
	return nil
}

// VerifyECDSAHeader verifies a signature using the EDCSA algorithm.
// It returns nil if the signature is valid, or an error if it is invalid.
func VerifyECDSAHeader(pubKey interface{}, tokens []string, signature []byte, sigType string) (bool, error) {
	parsedCryptoKey := pubKey.(ssh.CryptoPublicKey)
	pubCrypto := parsedCryptoKey.CryptoPublicKey()
	pkey, ok := pubCrypto.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("expected an ECDSA public key, got %T", pkey)
	}
	switch sigType {
	case "asn1":
		hashed := []byte(tokens[0])
		hashedSlice := hashed[:]
		verified := ecdsa.VerifyASN1(pkey, hashedSlice, signature)
		if verified != true {
			return false, fmt.Errorf("unable to verify ECDSA signature")
		}
	default:
		hashed := sha256.Sum256([]byte(tokens[0]))
		hashedSlice := hashed[:]

		var ecSig struct {
			R *big.Int
			S *big.Int
		}

		if err := ssh.Unmarshal(signature, &ecSig); err != nil {
			// We silently fail here, because this will happen whenever the client
			// did not send an ECDSA signature and we are trying to check it
			// against an ECDSA public key.
			return false, nil
		}

		verified := ecdsa.Verify(pkey, hashedSlice, ecSig.R, ecSig.S)
		if verified != true {
			return false, fmt.Errorf("unable to verify ECDSA signature")
		}
	}
	return true, nil
}
