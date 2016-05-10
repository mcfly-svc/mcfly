package api

import "github.com/mikec/msplapi/models"

type LoginReq struct {
	Token    string `json:"token" validate:"nonzero"`
	Provider string `json:"provider" validate:"nonzero"`
}

func (lr *LoginReq) AuthProvider() string {
	return lr.Provider
}

type LoginResp struct {
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

// Login with a provider access token
func (handlers *Handlers) Login(r *Responder, ctx *RequestContext) {

	loginReq := ctx.RequestData.(*LoginReq)
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
	u, err = handlers.db.GetUserByProviderToken(pt)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	if u == nil { // if user doesn't exist
		u = &models.User{
			Name:        td.UserName,
			AccessToken: handlers.generateToken(),
		}
		if err = handlers.db.SaveUser(u); err != nil {
			r.RespondWithServerError(err)
			return
		}
		if err = handlers.db.SetUserProviderToken(u.ID, pt); err != nil {
			r.RespondWithServerError(err)
			return
		}
	}

	r.Respond(&LoginResp{
		Name:        u.Name,
		AccessToken: u.AccessToken,
	})
}
