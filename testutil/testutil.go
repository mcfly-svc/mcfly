package testutil

import (
  "github.com/stretchr/testify/assert"

  "fmt"
  "encoding/json"
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

func (rt *ResponseTest) ExpectResponseBody(expectedBody interface{}) {
  b, err := ioutil.ReadAll(rt.Response.Body)
  if err != nil {
  	panic(err)
  }
  actualBody := string(b)

  var expBodyStr string
	switch v := expectedBody.(type) {
		case string:
			expBodyStr = v
		default:
			bodyBytes, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			expBodyStr = string(bodyBytes)
	}
	expBodyStr = fmt.Sprintf("%s\n", expBodyStr)

  assert.Equal(rt.Test, expBodyStr, actualBody)
}
