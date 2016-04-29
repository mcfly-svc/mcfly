package api_test

import (
	"testing"

	"github.com/mikec/marsupi-api/api"
)

func TestLoginInvalidJson(t *testing.T) {
	RunPostInvalidJsonTest(t, "login")
}

func TestLoginNoToken(t *testing.T) {
	RunMissingPostParamTest(t, "login", `{ "token_type":"jabroni.com" }`, "token")
}

func TestLoginNoTokenType(t *testing.T) {
	RunMissingPostParamTest(t, "login", `{ "token":"abc123" }`, "token_type")
}

func TestLoginUnsupportedTokenType(t *testing.T) {
	RunPostErrorTest(&ApiErrorTest{
		t,
		"login",
		`{ "token":"abc123", "token_type":"junk-service" }`,
		"an unsupported token type",
		"an unsupported token type error",
		api.NewUnsupportedTokenTypeErr("junk-service"),
	})
}

func TestLoginInvalidToken(t *testing.T) {
	RunPostErrorTest(&ApiErrorTest{
		t,
		"login",
		`{ "token":"badtoken", "token_type":"jabroni.com" }`,
		"an invalid token",
		"an invalid token error",
		api.NewInvalidTokenErr("jabroni.com"),
	})
}
