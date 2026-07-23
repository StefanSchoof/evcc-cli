package cmd

import (
	"fmt"

	"evcc-cli/internal/buildinfo"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "version=%s commit=%s date=%s\n", buildinfo.Version, buildinfo.Commit, buildinfo.Date)
		},
	}
}
