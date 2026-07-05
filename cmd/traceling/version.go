package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set at build time via -ldflags (see Makefile).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("traceling %s (commit %s, built %s)\n", version, commit, date)
		},
	}
}
