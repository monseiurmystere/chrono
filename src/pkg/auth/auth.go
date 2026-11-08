package auth

// AuthState represents the current state of the PAM authentication worker.
type AuthState string

const (
	AuthIdle           AuthState = "idle"
	AuthPrompting      AuthState = "prompting"
	AuthAuthenticating AuthState = "authenticating"
	AuthSuccess        AuthState = "success"
	AuthFailure        AuthState = "failure"
	AuthLockedOut      AuthState = "locked_out"
)

// AuthEvent represents an authentication progress or result notification.
type AuthEvent struct {
	State   AuthState
	Message string
	Prompt  string
}

// Authenticator defines the interface for interacting with the authentication engine.
type Authenticator interface {
	Start(username string) (<-chan AuthEvent, error)
	SupplySecret(secret []byte)
	Cancel()
}

// ClearSecret securely zeroes out memory buffers containing sensitive credentials.
func ClearSecret(secret []byte) {
	for i := range secret {
		secret[i] = 0
	}
}
