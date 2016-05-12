package provider

import "fmt"

type Dropbox struct {
	*SourceProviderConfig
}

func (self *Dropbox) Key() string {
	return "dropbox"
}

func (self *Dropbox) GetTokenData(string) (*TokenDataResponse, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (self *Dropbox) GetProjectData(token string, projectHandle string) (*ProjectData, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (self *Dropbox) GetProjects(token string, username string) ([]ProjectData, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (self *Dropbox) CreateProjectUpdateHook(token string, projectHandle string) error {
	return fmt.Errorf("Not Implemented")
}
