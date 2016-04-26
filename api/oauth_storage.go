package api

import(
	"fmt"

  "github.com/RangelReale/osin"

	"github.com/mikec/marsupi-api/models"
)

type OAuthStorage struct {
	db *models.DB
}

/*type OAuthData struct {
	User *model.User
}*/

func (s *OAuthStorage) GetClient(id string) (osin.Client, error) {
	return &osin.DefaultClient{
		Id:          "MCLOVIN",
		Secret:      "abc",
		RedirectUri: "http://localhost:8080/junk",
	}, nil
	/*secret, ok := appClients[id]
	if !ok {
		return nil, fmt.Errorf("Client not found")
	}

	return &osin.DefaultClient{
		Id:          id,
		Secret:      secret,
		RedirectUri: "about:blank",
	}, nil*/
}

func (s *OAuthStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	return fmt.Errorf("SaveAuthorize not implemented")
}

func (s *OAuthStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	return nil, fmt.Errorf("LoadAuthorize not implemented")
}

func (s *OAuthStorage) RemoveAuthorize(code string) error {
	return fmt.Errorf("RemoveAuthorize not implemented")
}

func (s *OAuthStorage) SaveAccess(data *osin.AccessData) error {
	fmt.Printf("DATA:%+v",data)

	return fmt.Errorf("SaveAccess not implemented")
	/*t := &model.AccessToken{
		Token:    data.AccessToken,
		ClientID: data.Client.GetId(),
	}

	if data.UserData != nil {
		d := data.UserData.(OAuthData)
		u := d.User

		if u != nil && !u.Deleted {
			t.UserID = u.ID
		}
	}

	err := s.db.InsertAccessToken(t)
	if err != nil {
		return err
	}

	return nil*/
}

func (s *OAuthStorage) LoadAccess(code string) (*osin.AccessData, error) {
	return nil, fmt.Errorf("LoadAccess not implemented")
}

func (s *OAuthStorage) RemoveAccess(code string) error {
	return fmt.Errorf("RemoveAccess not implemented")
}

func (s *OAuthStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	return nil, fmt.Errorf("LoadRefresh not implemented")
}

func (s *OAuthStorage) RemoveRefresh(token string) error {
	return fmt.Errorf("RemoveRefresh not implemented")
}

func (s *OAuthStorage) Clone() osin.Storage {
	return s
}

func (s *OAuthStorage) Close() {}
