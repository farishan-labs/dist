package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist/internal/metrics"
	"github.com/farishan-labs/dist/internal/registry"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check tools, auth, and registry health",
	RunE: func(cmd *cobra.Command, args []string) error {
		ok := true
		check := func(name string, fn func() error) {
			if err := fn(); err != nil {
				ok = false
				fmt.Printf("[FAIL] %s: %v\n", name, err)
			} else {
				fmt.Printf("[OK]   %s\n", name)
			}
		}

		required := []string{"gh", "gumroad", "cf"}
		optional := []string{"goreleaser", "wrangler", "npm"}
		for _, t := range required {
			tt := t
			check("tool:"+tt, func() error {
				if _, err := exec.LookPath(tt); err != nil {
					return err
				}
				return nil
			})
		}
		for _, t := range optional {
			tt := t
			if _, err := exec.LookPath(tt); err != nil {
				fmt.Printf("[WARN] tool:%s: not installed (optional)\n", tt)
			} else {
				fmt.Printf("[OK]   tool:%s\n", tt)
			}
		}

		check("registry", func() error {
			reg, err := registry.Load(registry.DefaultPath())
			if err != nil {
				return err
			}
			if err := reg.Validate(); err != nil {
				return err
			}
			if len(reg.Channels) == 0 {
				return fmt.Errorf("no channels configured")
			}
			return nil
		})

		check("gh-auth", func() error {
			c := exec.Command("gh", "auth", "status")
			if err := c.Run(); err != nil {
				return err
			}
			return nil
		})

		check("gumroad-auth", func() error {
			c := exec.Command("gumroad", "products", "list", "--json")
			if err := c.Run(); err != nil {
				return err
			}
			return nil
		})

		if os.Getenv("GUMROAD_ACCESS_TOKEN") == "" {
			fmt.Println("[WARN] env:GUMROAD_ACCESS_TOKEN: not set (optional; needed for views)")
		} else {
			fmt.Println("[OK]   env:GUMROAD_ACCESS_TOKEN")
		}

		if os.Getenv("GOATCOUNTER_URL") != "" {
			check("goatcounter-api", func() error {
				return metrics.VerifyGoatCounter()
			})
		} else {
			fmt.Println("[WARN] env:GOATCOUNTER_URL: not set (optional; needed for traffic)")
		}

		if !ok {
			fmt.Println("\nSome checks failed.")
			return fmt.Errorf("doctor: failures present")
		}
		fmt.Println("\nAll checks passed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
