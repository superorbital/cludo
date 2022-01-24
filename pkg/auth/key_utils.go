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

	pubkey, ok := pubCrypto.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expected an RSA public key, got %T", pubCrypto)
	}

	// Finally, we can convert back to an *rsa.PublicKey
	return pubkey, nil
}

func EncodeAuthorizedKey(key *rsa.PublicKey) (string, error) {
	pub, err := ssh.NewPublicKey(key)
	if err != nil {
		return "", err
	}

	return string(ssh.MarshalAuthorizedKey(pub)), nil
}

// getPassphrase will prompt the user for a passphrase
// It returns a []bytes containing the user input and any errors encountered.
func getPassphrase(path string) ([]byte, error) {
	fmt.Printf("Please enter the SSH key passphrase for %s: ", path)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return []byte(""), err
	}

	fmt.Printf("\n")
	password := string(bytePassword)
	return []byte(strings.TrimSpace(password)), nil
}

// DecodePrivateKey tries to decode the given private key.
// It wil try and handle passphrase-protected keys when encountered.
// It returns the decoded private key and any errors that were encountered.
func DecodePrivateKey(path string, encoded []byte, interactive bool) (*rsa.PrivateKey, error) {
	var parsedKey interface{}
	var err error

	// See if we can decode without a passphrase
	parsedKey, err = ssh.ParseRawPrivateKey(encoded)
	if err != nil && interactive == true {
		// In this case, let's try to prompt for the passphrase
		var passphrase []byte
		if utils.DetectTerminal() == true {
			passphrase, err = getPassphrase(path)
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
