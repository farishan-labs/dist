package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Print distribution metrics as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		// v0: stub metrics. Real implementation pulls Gumroad analytics + GoatCounter.
		m := map[string]any{
			"status":   "stub",
			"views":    0,
			"sales":    0,
			"revenue":  0,
			"note":     "Implement Gumroad get_analytics + GoatCounter API in v0.2",
		}
		b, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		printJSON(b)
		return nil
	},
}
