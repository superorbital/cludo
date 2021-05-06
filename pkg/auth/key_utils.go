package auth

import (
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func DecodeAuthorizedKey(encoded []byte) (*rsa.PublicKey, error) {
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(encoded)
	if err != nil {
		return nil, err
	}

	// To get back to an *rsa.PublicKey, we need to first upgrade to the ssh.CryptoPublicKey interface
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	return pubCrypto.(*rsa.PublicKey), nil
}

func EncodeAuthorizedKey(key *rsa.PublicKey) (string, error) {
	pub, err := ssh.NewPublicKey(key)
	if err != nil {
		return "", err
	}

	return string(ssh.MarshalAuthorizedKey(pub)), nil
}

func DecodePrivateKey(encoded []byte, password []byte) (*rsa.PrivateKey, error) {
	var parsedKey interface{}
	var err error
	if password != nil {
		parsedKey, err = ssh.ParseRawPrivateKeyWithPassphrase(encoded, password)
	} else {
		parsedKey, err = ssh.ParseRawPrivateKey(encoded)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %v, %#v", err, encoded)
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Parsed private key is not an RSA private key: %#v", parsedKey)
	}

	return privateKey, nil
}
