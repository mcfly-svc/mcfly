package provider

type TokenDataResponse struct {
	IsValid          bool
	Provider         string
	ProviderUsername string
	UserName         *string
}

// AuthProvider is a service used for authenticating to msplapi
type AuthProvider interface {
	Provider

	// GetTokenData given an auth token returns associated data
	GetTokenData(string) (*TokenDataResponse, error)
}

func GetAuthProviders() map[string]AuthProvider {
	github := GitHub{GitHubClient: &GoGitHubClient{}}
	authProviders := make(map[string]AuthProvider)
	authProviders[github.Key()] = &github
	return authProviders
}
