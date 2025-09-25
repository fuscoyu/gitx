package main

import (
	"fmt"
	"runtime"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// VersionInfo returns version information
func VersionInfo() string {
	return fmt.Sprintf("GitX %s\nBuild Time: %s\nGit Commit: %s\nGo Version: %s\nOS/Arch: %s/%s",
		version, buildTime, gitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// Version returns the version string
func Version() string {
	return version
}
