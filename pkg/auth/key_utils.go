package auth

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
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

func GetPassphrase() ([]byte, error) {
	fmt.Printf("\nUnable to decode SSH key on first attempt.\n")
	fmt.Print("Please enter the correct SSH passphrase: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return []byte(""), err
	}

	fmt.Printf("\n")
	password := string(bytePassword)
	return []byte(strings.TrimSpace(password)), nil
}

func DecodePrivateKey(encoded []byte, password []byte, interactive bool) (*rsa.PrivateKey, error) {
	var parsedKey interface{}
	var err error

	// See if we can decode without a passphrase
	parsedKey, err = ssh.ParseRawPrivateKey(encoded)
	if err != nil && len(password) > 0 {
		// if not, let's try the passphrase the user provided
		parsedKey, err = ssh.ParseRawPrivateKeyWithPassphrase(encoded, password)
	}
	// If we still have an error the passphrase is likely unset or wrong,
	// so let's prompt for it instead.
	if err != nil && interactive == true {
		var passphrase []byte
		passphrase, err = GetPassphrase()
		if err == nil {
			parsedKey, err = ssh.ParseRawPrivateKeyWithPassphrase(encoded, passphrase)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %v", err)
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Parsed private key is not an RSA private key: %#v", parsedKey)
	}

	return privateKey, nil
}
