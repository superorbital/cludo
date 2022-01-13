package cludo

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"

	"github.com/superorbital/cludo/client/environment"
	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/config"
)

const CludoProfileEnvVar = "CLUDO_PROFILE"

// GenerateEnvironment generates an environment bundle on a remote cludod service.
func GenerateEnvironment(cc *config.ClientConfig, target string, keys []string, debug bool, dryRun bool) (models.ModelsEnvironmentBundle, error) {
	cludodClient, err := config.NewClient(target, debug)
	if err != nil {
		return nil, err
	}
	signer, err := cc.NewDefaultSigner(keys)
	if err != nil {
		return nil, err
	}

	params := environment.NewGenerateEnvironmentParams().WithDefaults().WithBody(&models.ModelsEnvironmentRequest{
		Target: target,
	})

	if dryRun {
		log.Printf("Dry run enabled: Would send %#v", params)
		return nil, nil
	}

	response, err := cludodClient.Environment.GenerateEnvironment(params, signer.CludoAuth())
	if err != nil {
		return nil, fmt.Errorf("[1] Failed to generate environment: %v, %#v", err, response)
	}

	if response != nil && response.Payload != nil {
		return response.Payload.Bundle, nil
	}

	return nil, nil
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
