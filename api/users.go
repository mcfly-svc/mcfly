package api

/*

PROBABLY GARBAGE

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mikec/marsupi-api/logging"
	"github.com/mikec/marsupi-api/models"
)

// TODO: Get rid of this, replace with login endpoint
//
// curl -X POST http://localhost:8080/api/0/users -d '{"github_token":"xxxxxx"}'
func (handlers *Handlers) UserPost(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	var usrReq CreateUserReq
	err := DecodeRequest(req, &usrReq)
	if err != nil {
		logging.LogInternalError("UserPost", err)
		r.RespondWithError(NewInvalidJsonErr())
		return
	}

	if usrReq.GitHubToken == nil {
		r.RespondWithError(NewMissingParamErr("github_token"))
		return
	}
	ghToken := *usrReq.GitHubToken

	u, err := handlers.github.GetAuthenticatedUser(ghToken)
	if err != nil {
		logging.LogInternalError("UserPost", err)
		r.RespondWithError(NewInvalidGitHubTokenErr(ghToken))
		return
	}

	if err := handlers.db.SaveUser(u); err != nil {
		qErr := err.(*models.QueryExecError)
		switch qErr.Name {
		case "unique_violation":
			r.RespondWithError(NewDuplicateCreateErr("user"))
		default:
			logging.LogInternalError("UserPost", err)
			r.RespondWithError(NewCreateFailedErr("user"))
		}
		return
	}

	r.Respond(u)
}

// curl -X GET http://localhost:8080/api/0/users/1
func (handlers *Handlers) UserGet(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	vars := mux.Vars(req)
	user_id := vars["user_id"]
	id, err := strconv.ParseInt(user_id, 10, 64)

	if id == 0 || err != nil {
		r.RespondWithError(NewInvalidParamErr("ID", user_id))
		return
	}

	user, err := handlers.db.GetUserById(id)
	if err != nil {
		logging.LogInternalError("UserGet", err)
		r.RespondWithError(NewGetEntityErr("user", id))
		return
	}

	r.Respond(user)
}

// curl -X DELETE http://localhost:8080/api/0/users/1
func (handlers *Handlers) UserDelete(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	vars := mux.Vars(req)
	user_id := vars["user_id"]
	id, err := strconv.ParseInt(user_id, 10, 64)

	if err != nil {
		logging.LogInternalError("UserDelete", err)
		r.RespondWithError(NewInvalidParamErr("user_id", user_id))
		return
	}

	if err := handlers.db.DeleteUser(id); err != nil {
		logging.LogInternalError("UserDelete", err)
		r.RespondWithError(NewDeleteFailedErr("user"))
		return
	}

	r.RespondWithSuccess()
}
*/
