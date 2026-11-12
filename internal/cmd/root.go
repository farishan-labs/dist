package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "dist",
	Short:   "FHL distribution layer: ship and measure across channels",
	Long:    "dist is the single command surface for shipping and measuring across FHL distribution channels.",
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(shipCmd)
	rootCmd.AddCommand(metricsCmd)
}

func printJSON(b []byte) {
	fmt.Println(string(b))
}
