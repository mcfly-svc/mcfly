package api

import (
    "net/http"
    "log"

    "github.com/mikec/marsupi-api/models"
    "github.com/mikec/marsupi-api/github"
)

type CreateUserReq struct {
	GitHubToken 				*string				`json:"github_token"`
	// BitbucketToken 	*string				`json:"bitbucket_token"`
	// DropboxToken			*string				`json:"dropbox_token"`
}

// curl -X POST http://localhost:8080/api/0/users -d '{"github_token":"xxxxxx"}'
func (handlers *Handlers) UsersPost(w http.ResponseWriter, req *http.Request) {
	var usrReq CreateUserReq
	err := DecodeRequest(req, &usrReq); if err != nil {
		log.Println(err)
    writeErrorResponse(w, InvalidJsonApiErr)
    return
	}

	if usrReq.GitHubToken == nil {
		msg := "Missing github_token parameter"
		log.Println(msg)
		writeErrorResponse(w, msg)
		return
	}

	client := github.GetAuthClient(*usrReq.GitHubToken)

  user, _, err := client.Users.Get("")
  if err != nil {
  	log.Println(err)
  	writeErrorResponse(w, "Get user from GitHub failed")
  	return
  }

  mUser := &models.User{
  	Name: *user.Name,
  	GitHubUsername: *user.Login,
  	GitHubToken: *usrReq.GitHubToken,
  }

  newUser, err := handlers.db.SaveUser(mUser)
  if err != nil {
  	log.Println(err)
  	writeErrorResponse(w, "Create new user failed")
  	return
  }

	writeResponse(w, newUser)
}

// curl -X GET http://localhost:8080/api/0/users
func (handlers *Handlers) UsersGet(w http.ResponseWriter, req *http.Request) {
    users, err := handlers.db.GetUsers()
    if err != nil {
        log.Println(err)
        writeErrorResponse(w, "Failed to get users")
        return
    }

    writeResponse(w, users)
}


