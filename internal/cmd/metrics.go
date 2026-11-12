package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist/internal/metrics"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Print distribution metrics as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		rep := metrics.Collect()
		b, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			return err
		}
		printJSON(b)
		return nil
	},
}
