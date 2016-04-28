package api_test

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLoginNoToken(t *testing.T) {
	res, err := apiClient.Context.DoPost("login", nil, nil)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("LOGIN RESP:%+v\n", string(b))
}
