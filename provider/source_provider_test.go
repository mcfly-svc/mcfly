package provider_test

import (
	"testing"

	"github.com/mcfly-svc/mcfly/provider"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectUpdateHookUrl(t *testing.T) {
	s := provider.SourceProviderConfig{"api/1/{provider}/something", "mock_webhook_secret"}
	url := s.GetProjectUpdateHookUrl("jabroni.com")
	assert.Equal(t, "api/1/jabroni.com/something", url)
	// TODO: test for WebhookSecret
}
