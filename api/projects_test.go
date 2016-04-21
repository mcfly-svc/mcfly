package api_test

import (
	"testing"
)

var tests EndpointTests

func init() {
  tests = EndpointTests{
    apiClient.Projects,
    `{"service": "github", "username": "mikec", "name": "example-1"}`,
    `{"service": "github", "username": "mikec", "name": "example-2"}`,
  }
}

func TestEndpointOperations(t *testing.T) {
	tests.RunAll(t)
}
