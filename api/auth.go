package api

import (
	"github.com/mcfly-svc/mcfly/api/apidata"
	"github.com/mcfly-svc/mcfly/models"
)

// Login with a provider access token
func (handlers *Handlers) Login(r *Responder, ctx *RequestContext) {

	loginReq := ctx.RequestData.(*apidata.LoginReq)
	authProvider := *ctx.AuthProvider

	td, err := authProvider.GetTokenData(loginReq.Token)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	if !td.IsValid {
		r.RespondWithError(NewInvalidTokenErr(loginReq.Provider))
		return
	}

	pt := &models.ProviderAccessToken{
		Provider:         td.Provider,
		ProviderUsername: td.ProviderUsername,
		AccessToken:      loginReq.Token,
	}

	var u *models.User
	u, err = handlers.DB.GetUserByProviderToken(pt)
	if err != nil && err != models.ErrNotFound {
		r.RespondWithServerError(err)
		return
	}
	if u == nil { // if user doesn't exist
		u = &models.User{
			Name:        td.UserName,
			AccessToken: handlers.GenerateToken(),
		}
		if err = handlers.DB.SaveUser(u); err != nil {
			r.RespondWithServerError(err)
			return
		}
		if err = handlers.DB.SetUserProviderToken(u.ID, pt); err != nil {
			r.RespondWithServerError(err)
			return
		}
	}

	r.Respond(&apidata.LoginResp{
		Name:        u.Name,
		AccessToken: u.AccessToken,
	})
}
