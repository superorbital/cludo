package config

import (
	"testing"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/stretchr/testify/assert"
)

func TestValidClientGenerate(t *testing.T) {
	newClient, err := NewClient("https://valid.hostname.com/profile", false)
	assert.NoError(t, err, "NewClient should not error when generated")
	transportRuntime := newClient.Transport.(*httptransport.Runtime)
	assert.Equal(t, "valid.hostname.com", transportRuntime.Host, "Transport host should be set")
}
