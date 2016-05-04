package api

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/mikec/msplapi/models"
	"gopkg.in/validator.v2"
)

type LoginReq struct {
	Token     string `json:"token" validate:"nonzero"`
	TokenType string `json:"token_type" validate:"nonzero"`
}

type LoginResp struct {
	AccessToken string `json:access_token`
}

// Access token endpoint
func (handlers *Handlers) Login(w http.ResponseWriter, req *http.Request) {

	r := Responder{w}

	var loginReq LoginReq
	err := DecodeRequest(req, &loginReq)
	if err != nil {
		r.RespondWithError(NewInvalidJsonErr())
		return
	}

	if err := validator.Validate(loginReq); err != nil {
		errs := err.(validator.ErrorMap)
		var badParam string
		if len(errs["Token"]) > 0 {
			badParam = "token"
		} else if len(errs["TokenType"]) > 0 {
			badParam = "token_type"
		}
		r.RespondWithError(NewMissingParamErr(badParam))
		return
	}

	authProvider := handlers.authProviders[loginReq.TokenType]
	if authProvider == nil {
		r.RespondWithError(NewUnsupportedTokenTypeErr(loginReq.TokenType))
		return
	}

	td, err := authProvider.GetTokenData(loginReq.Token)
	if err != nil {
		r.RespondWithUnknownError("Login: authProvider.GetTokenData", err)
		return
	}

	if !td.IsValid {
		r.RespondWithError(NewInvalidTokenErr(loginReq.TokenType))
		return
	}

	pt := &models.ProviderAccessToken{
		td.Provider,
		td.ProviderUsername,
		loginReq.Token,
	}

	var u *models.User
	u, err = handlers.db.GetUserByProviderToken(pt)
	if err != nil {
		r.RespondWithUnknownError("Login: GetUserByProviderToken", err)
		return
	}
	if u == nil { // if user doesn't exist
		u = &models.User{
			Name:        td.UserName,
			AccessToken: generateAccessToken(),
		}
		if err = handlers.db.SaveUser(u); err != nil {
			r.RespondWithUnknownError("Login: SaveUser", err)
			return
		}
		if err = handlers.db.SetUserProviderToken(u.ID, pt); err != nil {
			r.RespondWithUnknownError("Login: SetUserProviderToken", err)
			return
		}
	}

	r.Respond(&LoginResp{u.AccessToken})
}

func generateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
