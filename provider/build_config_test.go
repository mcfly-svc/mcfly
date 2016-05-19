package provider_test

import (
	"testing"

	"github.com/mikec/msplapi/provider"
	"github.com/stretchr/testify/assert"
)

func TestNewBuildConfig(t *testing.T) {
	tests := []struct {
		JSON           []byte
		ExpectConfig   *provider.BuildConfig
		ExpectWarnings []provider.BuildConfigWarning
	}{
		{
			JSON: []byte(`{ "site":"/public" }`),
			ExpectConfig: &provider.BuildConfig{
				JSON: []byte(`{ "site":"/public" }`),
				Site: "/public",
			},
			ExpectWarnings: nil,
		},
		{
			JSON:         nil,
			ExpectConfig: provider.NewDefaultBuildConfig(),
			ExpectWarnings: []provider.BuildConfigWarning{
				{Message: "missing config file"},
			},
		},
		{
			JSON: []byte(`{jnk_json}`),
			ExpectConfig: &provider.BuildConfig{
				JSON: []byte(`{jnk_json}`),
				Site: "/",
			},
			ExpectWarnings: []provider.BuildConfigWarning{
				{
					Message: "invalid character 'j' looking for beginning of object key string",
					Line:    intPtr(1),
					Char:    intPtr(1),
				},
			},
		},
	}

	for _, test := range tests {
		bc, warnings := provider.NewBuildConfig(test.JSON)
		assert.Equal(t, test.ExpectConfig, bc)
		assert.Equal(t, test.ExpectWarnings, warnings)
	}
}
