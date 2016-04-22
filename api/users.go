package api

import (
    "net/http"
    "log"
    "fmt"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/mikec/marsupi-api/models"
)

type CreateUserReq struct {
	GitHubToken 				*string				`json:"github_token"`
	// BitbucketToken 	*string				`json:"bitbucket_token"`
	// DropboxToken			*string				`json:"dropbox_token"`
}

// curl -X POST http://localhost:8080/api/0/users -d '{"github_token":"xxxxxx"}'
func (handlers *Handlers) UserPost(w http.ResponseWriter, req *http.Request) {
  r := Responder{w}

	var usrReq CreateUserReq
	err := DecodeRequest(req, &usrReq); if err != nil {
		log.Println(err)
    r.RespondWithError(InvalidJsonApiErr)
    return
	}

	if usrReq.GitHubToken == nil {
		msg := "Missing github_token parameter"
		log.Println(msg)
		r.RespondWithError(msg)
		return
	}

  u, err := handlers.github.GetAuthenticatedUser(*usrReq.GitHubToken)
  if err != nil {
    log.Println(err)
    r.RespondWithError("github.GetAuthenticatedUser failed")
    return
  }

  if err := handlers.db.SaveUser(u); err != nil {
    qErr := err.(*models.QueryExecError)
    switch qErr.Name {
    case "unique_violation":
      r.RespondWithError("User already exists")
    default:
      r.RespondWithError("Failed to save user")
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
    apiErr := &ApiError{}
    apiErr.InvalidParam("ID", user_id)
    r.RespondWithError(*apiErr)
    return
  }
  
  user, err := handlers.db.GetUserById(id)
  if err != nil {
    log.Println(err)
    r.RespondWithError(fmt.Sprintf("Failed to get user where id=%d", id))
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
        log.Println(err)
        r.RespondWithError(fmt.Sprintf("%s is not a valid User ID", user_id))
        return
    }

    if err := handlers.db.DeleteUser(id); err != nil {
        log.Println(err)
        r.RespondWithError("Failed to delete User")
        return
    }

    r.RespondWithSuccess()
}

