package api_test

import (
	"testing"
)

var testRunner EndpointTestRunner

func init() {
  testRunner = EndpointTestRunner{
    apiClient.Projects,
    `{"service": "github", "username": "mikec", "name": "example-1"}`,
    `{"service": "github", "username": "mikec", "name": "example-2"}`,
  }
}

func TestEndpointOperations(t *testing.T) {
	testRunner.RunAll(t)
}
