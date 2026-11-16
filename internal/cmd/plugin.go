package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist/internal/plugins"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage dist metric plugins",
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := plugins.Load()
		if err != nil {
			return err
		}
		for _, p := range m.Plugins {
			fmt.Printf("%-14s %-12s enabled=%-5v binary=%s args=%v\n", p.Name, "plugin", p.Enabled, p.Binary, p.Args)
		}
		return nil
	},
}

var pluginAddCmd = &cobra.Command{
	Use:   "add <name> <binary> [args...]",
	Short: "Add a plugin",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		binary := args[1]
		var a []string
		if len(args) > 2 {
			a = args[2:]
		}
		if err := plugins.Add(name, binary, a); err != nil {
			return err
		}
		fmt.Printf("added plugin %s\n", name)
		return nil
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := plugins.Remove(args[0]); err != nil {
			return err
		}
		fmt.Printf("removed plugin %s\n", args[0])
		return nil
	},
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enable a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return plugins.SetEnabled(args[0], true)
	},
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disable a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return plugins.SetEnabled(args[0], false)
	},
}

var pluginRunCmd = &cobra.Command{
	Use:   "run <name>",
	Short: "Run a plugin and print its output",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := plugins.Load()
		if err != nil {
			return err
		}
		p, idx := plugins.Find(m, args[0])
		if idx < 0 {
			return fmt.Errorf("plugin %s not found", args[0])
		}
		out, err := plugins.Run(p)
		if err != nil {
			return err
		}
		var raw map[string]any
		if err := json.Unmarshal(out, &raw); err == nil {
			b, _ := json.MarshalIndent(raw, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Println(string(out))
		}
		return nil
	},
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginAddCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)
	pluginCmd.AddCommand(pluginEnableCmd)
	pluginCmd.AddCommand(pluginDisableCmd)
	pluginCmd.AddCommand(pluginRunCmd)
	rootCmd.AddCommand(pluginCmd)
}
