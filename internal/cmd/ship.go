package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist/internal/registry"
	"github.com/farishan-labs/dist/internal/ship"
)

var (
	shipDryRun  bool
	shipYes     bool
	shipChannel []string
	shipTag     string
	shipDir     string
	shipFile    string
)

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Orchestrate shipping across channels (dry-run by default)",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.Load(registry.DefaultPath())
		if err != nil {
			return err
		}
		if verr := reg.Validate(); verr != nil {
			return fmt.Errorf("registry invalid: %w", verr)
		}

		filter := map[string]bool{}
		for _, c := range shipChannel {
			filter[c] = true
		}
		channels := reg.Enabled()
		if len(filter) > 0 {
			var sel []registry.Channel
			for _, c := range channels {
				if filter[c.Name] {
					sel = append(sel, c)
				}
			}
			channels = sel
		}
		if len(channels) == 0 {
			fmt.Println("No channels to ship.")
			return nil
		}

		fmt.Println("Ship plan:")
		for _, c := range channels {
			fmt.Printf("  - %s (%s) dry_run=%v\n", c.Name, c.Type, shipDryRun)
		}

		if !shipDryRun && !shipYes {
			fmt.Print("Proceed? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(strings.ToLower(line))
			if line != "y" && line != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		results := ship.Run(channels, shipDryRun, ship.Options{Tag: shipTag, Dir: shipDir, File: shipFile})
		for _, r := range results {
			status := "OK"
			if r.Err != nil {
				status = "FAILED: " + r.Err.Error()
			} else if r.Skipped {
				status = "SKIPPED"
			}
			fmt.Printf("[%s] %s: %s\n", r.Channel, status, r.Summary)
		}

		if err := ship.AppendEvents(results); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not write event log: %v\n", err)
		}
		return nil
	},
}

func init() {
	shipCmd.Flags().BoolVar(&shipDryRun, "dry-run", true, "print planned actions without executing")
	shipCmd.Flags().BoolVar(&shipYes, "yes", false, "skip approval prompt")
	shipCmd.Flags().StringSliceVar(&shipChannel, "channel", nil, "limit to named channel(s)")
	shipCmd.Flags().StringVar(&shipTag, "tag", "", "release tag (github)")
	shipCmd.Flags().StringVar(&shipDir, "dir", "", "build/output dir override")
	shipCmd.Flags().StringVar(&shipFile, "file", "", "file to attach (gumroad)")
}
