package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist/internal/registry"
)

var statusJSON bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show channel/product inventory and last-shipped state",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.Load(registry.DefaultPath())
		if err != nil {
			return err
		}
		if verr := reg.Validate(); verr != nil {
			return fmt.Errorf("registry invalid: %w", verr)
		}
		if statusJSON {
			b, err := json.MarshalIndent(reg, "", "  ")
			if err != nil {
				return err
			}
			printJSON(b)
			return nil
		}
		if len(reg.Channels) == 0 {
			fmt.Println("No channels configured. Create", registry.DefaultPath())
			return nil
		}
		fmt.Printf("schema_version=%d channels=%d enabled=%d\n", reg.SchemaVersion, len(reg.Channels), len(reg.Enabled()))
		fmt.Printf("%-20s %-12s %-8s %-28s %s\n", "NAME", "TYPE", "ENABLED", "TARGET", "LAST_SHIPPED")
		for _, c := range reg.Channels {
			target := c.Repo
			if target == "" {
				target = c.ProductID
			}
			if target == "" {
				target = c.Project
			}
			if target == "" {
				target = c.Package
			}
			if target == "" {
				target = c.Platform
			}
			ls := c.LastShipped
			if ls == "" {
				ls = "-"
			}
			fmt.Printf("%-20s %-12s %-8v %-28s %s\n", c.Name, c.Type, c.Enabled, target, ls)
		}
		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "output as JSON")
}
