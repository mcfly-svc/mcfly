package api

import (
    "fmt"
    "net/http"
    "github.com/mikec/marsupi-api/models"
)

type LoginReq struct {
    GitHubToken                 *string             `json:"github_token"`
    // BitbucketToken   *string             `json:"bitbucket_token"`
    // DropboxToken         *string             `json:"dropbox_token"`
}

// Access token endpoint
func (handlers *Handlers) Login(w http.ResponseWriter, req *http.Request) {

    r := Responder{w}

    osinResp := handlers.osinServer.NewResponse()
    defer osinResp.Close()

    if ar := handlers.osinServer.HandleAccessRequest(osinResp, req); ar != nil {

        var loginReq LoginReq
        err := DecodeRequest(req, &loginReq); if err != nil {
            fmt.Println(err)
            r.RespondWithError(NewInvalidJsonErr())
            return
        }
        ghToken := loginReq.GitHubToken

        if ghToken == nil {
            r.RespondWithError(NewMissingParamErr("github_token"))
            return
        }

        u, err := handlers.github.GetAuthenticatedUser(*ghToken)
        if err != nil {
            fmt.Println(err)
            r.RespondWithError(NewInvalidGitHubTokenErr(*ghToken))
            return
        }

        savedUser, err := handlers.db.GetUserByGitHubToken(u.GitHubToken)
        if err != nil {
            qErr := err.(*models.QueryExecError)
            /*switch qErr.Name {
                case 
            }*/
            fmt.Println("USER ERROR", qErr.Name)
        }

        // SAVE USER, OR ADD TOKEN TO USER?? HERE OR IN oauth_storage??
        // and why is TestGetUser failing?
        ar.UserData = savedUser

        ar.Authorized = true
        handlers.osinServer.FinishAccessRequest(osinResp, req, ar)
    }

    if osinResp.IsError {
        if osinResp.InternalError != nil {
            fmt.Println("osin error:", osinResp.InternalError)
        }
        r.RespondWithError(NewApiErr(osinResp.StatusText))
        return
    }

    r.Respond(osinResp)

    //osin.OutputJSON(osinResp, w, req)
}
