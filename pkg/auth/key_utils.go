package auth

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/superorbital/cludo/pkg/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// ParseAuthorizedKey converts an authorized_key into a PublicKey
// It returns the PublicKey and any errors encountered.
func ParseAuthorizedKey(encoded []byte) (*ssh.PublicKey, error) {
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(encoded)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

// EncodeAuthorizedKey will encode the given public key
// into the authorized_key format
// It returns a string containing the authorized_key and any errors encountered.
func EncodeAuthorizedKey(pub *ssh.PublicKey) (string, error) {

	return string(ssh.MarshalAuthorizedKey(*pub)), nil
}

// getPassphrase will prompt the user for a passphrase
// It returns a []bytes containing the user input and any errors encountered.
func getPassphrase(path string) ([]byte, error) {
	fmt.Printf("Please enter the SSH key passphrase for %s: ", path)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return []byte(""), err
	}

	fmt.Printf("\n\n")
	password := string(bytePassword)
	return []byte(strings.TrimSpace(password)), nil
}

// DecodePrivateKey tries to decode the given private key.
// It will try and handle passphrase-protected keys when encountered.
// It returns the decoded private key and any errors that were encountered.
func DecodePrivateKey(path string, encoded []byte, interactive bool) (*interface{}, error) {
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

	return &parsedKey, nil
}
