package provider

import "fmt"

type Dropbox struct{}

func (self *Dropbox) Key() string {
	return "dropbox"
}

func (self *Dropbox) GetTokenData(string) (*TokenDataResponse, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (self *Dropbox) GetProjectData(token string, projectHandle string) (*ProjectData, error) {
	return nil, fmt.Errorf("Not Implemented")
}
