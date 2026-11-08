package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"io.tylr.chrono/pkg/providers"
)

var (
	// Version is populated during build via -ldflags
	Version   = "0.1.0-dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

func main() {
	configDir := flag.String("config", "", "Path to custom configuration directory")
	flag.StringVar(configDir, "c", "", "Path to custom configuration directory (shorthand)")

	testMode := flag.Bool("test", false, "Run in test mode without locking the Wayland session")
	flag.BoolVar(testMode, "t", false, "Run in test mode without locking the Wayland session (shorthand)")

	showVersion := flag.Bool("version", false, "Show version information and exit")
	flag.BoolVar(showVersion, "v", false, "Show version information and exit (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Chrono %s (commit: %s, built: %s)\n\n", Version, GitCommit, BuildDate)
		fmt.Fprintf(os.Stderr, "Usage: chrono [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("Chrono %s (commit: %s, built: %s)\n", Version, GitCommit, BuildDate)
		os.Exit(0)
	}

	fmt.Printf("Starting Chrono %s...\n", Version)

	ctx := providers.NewDefaultContext()
	fmt.Printf("Initialized default render context for user '%s' on '%s' (Time: %s)\n",
		ctx.User.Username, ctx.User.Hostname, ctx.Clock.Time)

	if *testMode {
		fmt.Println("[TEST MODE] Session lock will not be acquired.")
	}

	// Trap termination signals for clean shutdown in development
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Chrono initialized successfully. Awaiting events (press Ctrl+C to terminate)...")
	<-sigCh
	fmt.Println("\nShutting down Chrono.")
}
