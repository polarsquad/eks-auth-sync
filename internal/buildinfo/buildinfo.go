package buildinfo

import (
	"fmt"
	"io"
	"log"
)

/*
This file includes version information fed to the binary during build phase.
*/

// Version is the version string based on Git tags
var Version string

// GitHash is the Git commit hash present during build
var GitHash string

func PrintVersion(out io.Writer) {
	_, err := fmt.Fprintf(
		out,
		"%s / git:%s\n",
		withPlaceholder(Version),
		withPlaceholder(GitHash),
	)
	if err != nil {
		log.Print("failed to print version information: %w", err)
	}
}

func withPlaceholder(s string) string {
	if s == "" {
		return "<unknown>"
	}
	return s
}
