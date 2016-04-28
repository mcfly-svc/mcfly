package provider

import "fmt"

type Dropbox struct {}

func (self Dropbox) Key() string {
  return "dropbox"
}

func (self Dropbox) GetTokenData(token string) (*TokenDataResponse, error) {
  return nil, fmt.Errorf("Dropbox.GetTokenData Not Implemented")
}
