package api

import (
	"net/http"

	"github.com/mikec/msplapi/models"
)

type LoginReq struct {
	Token    string `json:"token" validate:"nonzero"`
	Provider string `json:"provider" validate:"nonzero"`
}

type LoginResp struct {
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

// Login with a provider access token
func (handlers *Handlers) Login(w http.ResponseWriter, req *http.Request) {

	r := Responder{w}

	var loginReq LoginReq
	err := DecodeRequest(req, &loginReq)
	if err != nil {
		r.RespondWithError(NewInvalidJsonErr())
		return
	}

	apiErr := ValidateRequestData(&loginReq)
	if apiErr != nil {
		r.RespondWithError(apiErr)
		return
	}

	authProvider := handlers.authProviders[loginReq.Provider]
	if authProvider == nil {
		r.RespondWithError(NewUnsupportedProviderErr(loginReq.Provider))
		return
	}

	td, err := authProvider.GetTokenData(loginReq.Token)
	if err != nil {
		r.RespondWithUnknownError("Login: authProvider.GetTokenData", err)
		return
	}

	if !td.IsValid {
		r.RespondWithError(NewInvalidTokenErr(loginReq.Provider))
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
			AccessToken: handlers.generateToken(),
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

	r.Respond(&LoginResp{
		Name:        u.Name,
		AccessToken: u.AccessToken,
	})
}
