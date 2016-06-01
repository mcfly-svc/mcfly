package provider_test

import (
	"testing"

	"github.com/mcfly-svc/mcfly/provider"
	"github.com/stretchr/testify/assert"
)

func TestNewBuildConfig(t *testing.T) {
	missingFile := provider.NewDefaultBuildConfig()
	missingFile.Warnings = []provider.BuildConfigWarning{
		{Message: "missing config file"},
	}

	tests := []struct {
		JSON         []byte
		ExpectConfig *provider.BuildConfig
	}{
		{
			JSON: []byte(`{ "site":"/public" }`),
			ExpectConfig: &provider.BuildConfig{
				JSON: []byte(`{ "site":"/public" }`),
				Properties: &provider.BuildConfigProperties{
					Site: "/public",
				},
				Warnings: make([]provider.BuildConfigWarning, 0),
			},
		},
		{
			JSON:         nil,
			ExpectConfig: missingFile,
		},
		{
			JSON: []byte(`{jnk_json}`),
			ExpectConfig: &provider.BuildConfig{
				JSON: []byte(`{jnk_json}`),
				Properties: &provider.BuildConfigProperties{
					Site: "/",
				},
				Warnings: []provider.BuildConfigWarning{
					{
						Message: "invalid character 'j' looking for beginning of object key string",
						Line:    intPtr(1),
						Char:    intPtr(1),
					},
				},
			},
		},
	}

	for _, test := range tests {
		bc := provider.NewBuildConfig(test.JSON)
		assert.Equal(t, test.ExpectConfig, bc)
	}
}
