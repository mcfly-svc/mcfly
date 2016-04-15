package testutil

import (
  "github.com/stretchr/testify/assert"

  "io/ioutil"
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

func (rt *ResponseTest) ExpectResponseBody(expectedBody string) {
  b, err := ioutil.ReadAll(rt.Response.Body)
  if err != nil {
  	panic(err)
  }
  actualBody := string(b)
  assert.Equal(rt.Test, expectedBody, actualBody)
}
