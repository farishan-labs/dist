package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist/internal/registry"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show channel/product inventory and last-shipped state",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.Load(registry.DefaultPath())
		if err != nil {
			return err
		}
		if len(reg.Channels) == 0 {
			fmt.Println("No channels configured. Create", registry.DefaultPath())
			return nil
		}
		fmt.Printf("%-20s %-12s %-8s %s\n", "NAME", "TYPE", "ENABLED", "TARGET")
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
			fmt.Printf("%-20s %-12s %-8v %s\n", c.Name, c.Type, c.Enabled, target)
		}
		return nil
	},
}
