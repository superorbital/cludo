package providers

import (
	"context"
	"testing"

	github "github.com/google/go-github/v42/github"
)

/*
We might want to consider using this library
for go-github testing at some point:
https://github.com/migueleliasweb/go-github-mock
*/

type mockClient struct {
	resp []*github.Key
}

func (m *mockClient) ListKeys(ctx context.Context, user string, opts *github.ListOptions) ([]*github.Key, *github.Response, error) {
	return m.resp, nil, nil
}

func TestListKeys(t *testing.T) {
	ctx := context.Background()
	mc := &mockClient{
		resp: []*github.Key{
			{
				ID:  github.Int64(1),
				Key: github.String("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAII2YNmNYja2pH/D3hr8wwFqtRXIAdCYA25kgQESiWoDc test-key-ed25519-nopp@example.com DO NOT USE!!"),
			},
			{
				ID:  github.Int64(2),
				Key: github.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDD2ll7CNP8qq5CqXda7vAGevhES4nFB9qLVE5dvZHw1Gsl0hBubRgvNpQskjUX2KSejGSxSK7gXdF0P5w1P84E4wofRM251C688lGLkLZRNl2THQUTHEgycSokEkCSXDzjDYqF6OTW6cQ7NdLvdhZFNWO+3nparZNpEsKoQ8Xu5jkYzQxUXbwc66psIZfD6yTT1QtvgcsgmGH3aL/3RaiznFwsVa5Zq4qA3wwgrnuV9Q/si4wJAuOVtnmAKgnxwrHbrb0WwPOhPa5P4QlntnWi6jexeJ39A9gKhwh7QMF00VSkiiWv4V4dsq6eXVT/lyNLy8p5yStkBV43l6DGHGvBMCZhJkNb85/M79S4KhjlHLeaVBLj648ImTrDqzrFD4G5k5cEL3VTsQ3kus839TH5xdOj8iCJ9x/lVCuSmrF53US6+dL4NXRSGixyo+DhUlWUlMhsor3cgmjA2bSS6h/h7c/qaM5Q7YZ7QZtGSMqA9sS+mV8YL5rd6XFckWUCyOjt6ehiDfGH6E49l4S6NxPdQI1OpSmdVVdlsspprSTGfoD/2voXFdMSLLFtudeglYjPteXoF+Nj80iOPUXRwHRgj2HKi+B7BS3c4TPL+MlrlTOznyq565rJ2DJ9u9BD0PB4Z7snOB7eUiI2OkKODQ6WWsH3u0qCO0kVjEotdB5yQQ== test-key-pp@example.com DO NOT USE!!"),
			},
		},
	}
	kc := &userKeyClient{ctx: ctx, userKeyClient: mc}
	response, err := kc.listKeys("test")
	if err != nil {
		t.Fatalf("listKeys: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("listKeys: expected 2 keys, got %d", len(response))
	}
}
