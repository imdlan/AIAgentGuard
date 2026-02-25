package cmd

import (
	"fmt"
	"os"

	"github.com/imdlan/AIAgentGuard/internal/policy"
	"github.com/spf13/cobra"
)

var (
	initConfig bool
	initPath   string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize AgentGuard configuration",
	Long: `Create a default policy configuration file in your home directory or current directory.
This configuration file defines security policies for sandbox execution.`,
	Example: `  agent-guard init
  agent-guard init --path .agent-guard.yaml`,
	RunE: initConfigFile,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&initConfig, "create", false, "Create configuration file")
	initCmd.Flags().StringVarP(&initPath, "path", "p", "", "Path for configuration file")
}

func initConfigFile(cmd *cobra.Command, args []string) error {
	// Determine config path
	configPath := initPath
	if configPath == "" {
		// Default to .agent-guard.yaml in current directory
		configPath = ".agent-guard.yaml"
	}

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration file already exists: %s", configPath)
	}

	// Get default configuration
	cfg := policy.GetDefaultConfig()

	// Save configuration
	if err := policy.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}

	fmt.Printf("✅ Configuration file created: %s\n", configPath)
	fmt.Println("\nYou can now customize the security policies in this file.")
	fmt.Println("Edit the file to adjust:")
	fmt.Println("  - Filesystem access rules")
	fmt.Println("  - Shell command permissions")
	fmt.Println("  - Network access control")
	fmt.Println("  - Environment variable blocking")
	fmt.Println("  - Sandbox settings")

	return nil
}
