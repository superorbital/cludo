package cludo

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/superorbital/cludo/client"
	"github.com/superorbital/cludo/client/environment"
	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/config"
)

const CludoProfileEnvVar = "CLUDO_PROFILE"

// CheckErr prints the msg with the prefix '[ERROR]' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintf(os.Stderr, "\n[ERROR] %s", msg)
		os.Exit(1)
	}
}

// GenerateEnvironment generates an environment bundle on a remote cludod service.
func GenerateEnvironment(cc *config.ClientConfig, target string, keys []string, debug bool, dryRun bool) (models.ModelsEnvironmentBundle, error) {
	cludodClient, err := config.NewClient(target, debug)
	if err != nil {
		return nil, err
	}

	var response *environment.GenerateEnvironmentOK

	for _, key := range keys {
		response, err := attemptKeyAuth(cludodClient, cc, target, key, debug, dryRun)
		if err != nil {
			continue
		}
		if response != nil && response.Payload != nil {
			return response.Payload.Bundle, nil
		}
	}

	if err != nil {
		return nil, fmt.Errorf("The server reported an error: %v, %#v", err, response)
	}

	return nil, fmt.Errorf("Did not find any authorized SSH keys. Ensure that your keys are unlocked in your SSH agent or talk to your administrator: %#v", response)
}

// attemptKeyAuth attempts to authenticate with the server using the given key.
// It returns the environment response and any error encountered.
func attemptKeyAuth(client *client.Cludod, cc *config.ClientConfig, target string, key string, debug bool, dryRun bool) (*environment.GenerateEnvironmentOK, error) {

	signer, err := cc.NewDefaultClientSigner(key)
	if err != nil {
		return nil, err
	}

	params := environment.NewGenerateEnvironmentParams().WithDefaults().WithBody(&models.ModelsEnvironmentRequest{
		Target: target,
	})

	if dryRun {
		fmt.Printf("[INFO] Dry run enabled: Would send %#v", params)
		return nil, nil
	}

	response, err := client.Environment.GenerateEnvironment(params, signer.CludoAuth())

	return response, err
}

func ProfileURL(targetURL string) (string, error) {
	if targetURL == "" {
		return "", fmt.Errorf("invalid target URL, empty")
	}
	u, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("Failed to parse target URL '%s': %v", targetURL, err)
	}
	return u.String(), nil
}

func ExecWithEnv(args []string, env models.ModelsEnvironmentBundle, inherit bool, target string) (int, error) {
	profileURL, err := ProfileURL(target)
	if err != nil {
		return -1, fmt.Errorf("Failed to generate profile url: %v", err)
	}

	cmd := exec.Command(args[0], args[1:]...)
	if inherit {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = []string{}
	}
	for k, v := range env {
		if v != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	if profileURL != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", CludoProfileEnvVar, profileURL))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}

	return 0, nil
}
