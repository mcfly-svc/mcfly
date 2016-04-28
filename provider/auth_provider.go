package provider

type TokenDataResponse struct {
	IsValid          bool
	Provider         string
	ProviderUsername string
	UserName         string
}

type AuthProvider interface {

	// get key identifier for this provider
	Key() string

	// get data from the provider based on a provider auth token
	GetTokenData(string) (*TokenDataResponse, error)
}
