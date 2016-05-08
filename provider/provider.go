package provider

// Provider represents an external service
type Provider interface {

	// Key returns an identifier string for this provider
	Key() string
}
