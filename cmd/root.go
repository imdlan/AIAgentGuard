package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent-guard",
	Short: "AI Agent Security Scanner and Sandbox",
	Long: `AgentGuard is a security tool for AI agents, CLI tools, and MCP servers.
It scans for permission risks, evaluates security threats, and provides
sandboxed execution environments.`,
}

var (
	// Global flags
	configFile string
	verbose    bool
	jsonOutput bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to policy configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "JSON output format")
}
