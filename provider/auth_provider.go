package provider

type TokenDataResponse struct {
	IsValid          bool
	Provider         string
	ProviderUsername string
	UserName         string
}

// AuthProvider is a service used for authenticating to msplapi
type AuthProvider interface {
	Provider

	// GetTokenData given an auth token returns associated data
	GetTokenData(string) (*TokenDataResponse, error)
}
