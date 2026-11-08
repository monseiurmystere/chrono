package providers

import "time"

// RenderContext holds the reactive state injected into DOM templates.
type RenderContext struct {
	User    UserContext
	Clock   ClockContext
	Battery BatteryContext
	Auth    AuthContext
}

// UserContext provides system user and machine identification.
type UserContext struct {
	Username    string
	Hostname    string
	DisplayName string
}

// ClockContext provides formatted time and date strings.
type ClockContext struct {
	Time      string
	Date      string
	Timestamp time.Time
}

// BatteryContext provides power metrics from Linux sysfs.
type BatteryContext struct {
	Present   bool
	Percent   int
	State     string // "charging", "discharging", "full", "unknown"
	Icon      string
	PowerDraw float64 // Watts
}

// AuthContext provides status updates and prompts from the PAM authentication worker.
type AuthContext struct {
	Status   string // "idle", "authenticating", "error", "locked_out", "success"
	Message  string
	Prompt   string
	Attempts int
}

// NewDefaultContext creates a default render context.
func NewDefaultContext() RenderContext {
	return RenderContext{
		User: UserContext{
			Username:    "user",
			Hostname:    "localhost",
			DisplayName: "User",
		},
		Clock: ClockContext{
			Time:      time.Now().Format("15:04"),
			Date:      time.Now().Format("Monday, January 2"),
			Timestamp: time.Now(),
		},
		Battery: BatteryContext{
			Present: false,
			Percent: 100,
			State:   "full",
			Icon:    "󰁹",
		},
		Auth: AuthContext{
			Status:   "idle",
			Message:  "",
			Prompt:   "Password:",
			Attempts: 0,
		},
	}
}
