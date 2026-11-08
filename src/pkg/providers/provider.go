package providers

import "context"

// ProviderUpdate contains an update event emitted by a state provider.
type ProviderUpdate struct {
	ProviderName string
	Data         any
}

// Provider defines the lifecycle interface for background state providers.
type Provider interface {
	Name() string
	Start(ctx context.Context, updateCh chan<- ProviderUpdate) error
	Stop() error
}
