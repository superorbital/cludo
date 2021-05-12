package cludo

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/superorbital/cludo/client/environment"
	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/config"
)

// GenerateEnvironment generates an environment bundle on a remote cludod service.
func GenerateEnvironment(cc *config.ClientConfig, debug bool, dryRun bool) (models.ModelsEnvironmentBundle, error) {
	cludodClient, err := cc.NewClient(debug)
	if err != nil {
		return nil, err
	}
	signer, err := cc.NewDefaultSigner()
	if err != nil {
		return nil, err
	}

	params := environment.NewGenerateEnvironmentParams().WithDefaults().WithBody(&models.ModelsEnvironmentRequest{})

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

func ExecWithEnv(args []string, env models.ModelsEnvironmentBundle, inherit bool) (int, error) {
	cmd := exec.Command(args[0], args[1:]...)
	if inherit {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = []string{}
	}
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}

	return 0, nil
}
