package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dryRun bool

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Orchestrate shipping across channels (dry-run by default)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			fmt.Println("[dry-run] would: goreleaser release --clean")
			fmt.Println("[dry-run] would: wrangler pages deploy")
			fmt.Println("[dry-run] would: sync Gumroad product metadata")
			fmt.Println("Pass --no-dry-run (not yet implemented) to execute; v0 is dry-run only.")
			return nil
		}
		fmt.Println("Execution not implemented in v0; use --dry-run (default).")
		return nil
	},
}

func init() {
	shipCmd.Flags().BoolVar(&dryRun, "dry-run", true, "print planned actions without executing")
}
