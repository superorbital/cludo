package auth

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"syscall"

	"github.com/superorbital/cludo/pkg/utils"
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
	fmt.Print("Please enter the SSH key passphrase: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return []byte(""), err
	}

	fmt.Printf("\n")
	password := string(bytePassword)
	return []byte(strings.TrimSpace(password)), nil
}

func DecodePrivateKey(encoded []byte, interactive bool) (*rsa.PrivateKey, error) {
	var parsedKey interface{}
	var err error

	// See if we can decode without a passphrase
	parsedKey, err = ssh.ParseRawPrivateKey(encoded)
	if err != nil && interactive == true {
		// In this case, let's try to prompt for the passphrase
		var passphrase []byte
		if utils.DetectTerminal() == true {
			passphrase, err = GetPassphrase()
			if err == nil {
				parsedKey, err = ssh.ParseRawPrivateKeyWithPassphrase(encoded, passphrase)
			}
		} else {
			fmt.Println("[WARN] No terminal detected. Skipping passphrase prompt.")
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
