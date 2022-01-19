package providers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/v42/github"
)

func CollectGithubPublicKeys(api_endpoint string, username string) (pubkeys []string, err error) {

	client := github.NewClient(nil)

	if api_endpoint != "" {
		client.BaseURL, err = url.Parse(api_endpoint)
		if err != nil {
			return nil, fmt.Errorf("Received error while trying to parse the Github API Endpoint from the server config: %s", api_endpoint)
		}
	}

	// list all keys for the user
	response, _, err := client.Users.ListKeys(context.Background(), username, nil)
	if err != nil {
		return nil, fmt.Errorf("Received error while trying to retrieve user (%s) public keys from Github: %v", username, err)
	}

	keys := []string{}
	for _, v := range response {
		keys = append(keys, v.GetKey())
	}

	return keys, nil
}
