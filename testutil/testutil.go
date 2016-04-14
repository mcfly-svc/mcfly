package testutil

import (
  "github.com/stretchr/testify/assert"

	"testing"
	"net/http"
)

type ResponseTest struct {
	Test 			*testing.T
	Response 	*http.Response
}

func (rt *ResponseTest) ExpectHttpStatus(code int) {
  assert.Equal(rt.Test, code, rt.Response.StatusCode, "Wrong HTTP status code in response")
}
