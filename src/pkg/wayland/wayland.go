package wayland

// Output represents a physical or virtual Wayland display output.
type Output struct {
	ID    uint32
	Name  string
	Width int
	Height int
	Scale int
}

// Surface represents a lock screen surface mapped to an output.
type Surface struct {
	Output *Output
	Width  int
	Height int
	Buffer []byte
}

// SessionLocker defines the Wayland ext-session-lock-v1 client interface.
type SessionLocker interface {
	Connect(display string) error
	Lock() error
	Unlock() error
	Outputs() []*Output
	CommitSurface(surface *Surface) error
	Close() error
}
