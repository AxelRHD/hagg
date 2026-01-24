// Package version provides application version information.
// The Version variable is set at build time via ldflags.
package version

// Version is the application version, set at build time.
// Format: tag (v1.0.0), tag-hash (v1.0.0-abc1234), or hash (abc1234)
// Append -dirty if there are uncommitted changes.
var Version = "dev"
