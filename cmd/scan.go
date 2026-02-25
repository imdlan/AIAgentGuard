package cmd

import (
	"fmt"

	"github.com/imdlan/AIAgentGuard/internal/report"
	"github.com/imdlan/AIAgentGuard/internal/risk"
	"github.com/imdlan/AIAgentGuard/internal/scanner"
	"github.com/imdlan/AIAgentGuard/pkg/model"
	"github.com/spf13/cobra"
)

var (
	scanCategory  string
	scanDirectory string
	pluginsOnly   bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Scan for security risks and permissions",
	Long: `Scan the current environment or a specific target for security risks.
Supports scanning:
  - Local environment permissions (filesystem, shell, network, secrets)
  - Plugin directories for risky code
  - Specific tools or agents`,
	Example: `  agent-guard scan
  agent-guard scan plugins --dir ./plugins
  agent-guard scan filesystem --category filesystem`,
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&scanCategory, "category", "t", "", "Scan category: filesystem, shell, network, secrets, plugins")
	scanCmd.Flags().StringVarP(&scanDirectory, "dir", "d", ".", "Directory to scan (for plugins)")
	scanCmd.Flags().BoolVarP(&pluginsOnly, "plugins", "p", false, "Scan plugins in the specified directory")
}

func runScan(cmd *cobra.Command, args []string) error {
	if pluginsOnly || scanCategory == "plugins" {
		return scanPlugins()
	}

	// Regular permission scan
	var scanResult interface{}

	if scanCategory != "" {
		// Scan specific category
		result := scanner.RunSpecificScan(scanCategory)
		scanResult = map[string]interface{}{
			"category": scanCategory,
			"risk":     result,
		}
	} else {
		// Scan all categories
		result := scanner.RunAllScans()
		analysis := risk.Analyze(result)
		scanResult = analysis
	}

	// Output results
	if jsonOutput {
		switch v := scanResult.(type) {
		case model.ScanReport:
			return report.PrintJSON(v)
		default:
			fmt.Printf("%#v\n", scanResult)
		}
	} else {
		switch v := scanResult.(type) {
		case model.ScanReport:
			report.PrintConsole(v)
		default:
			fmt.Printf("Scan result: %v\n", scanResult)
		}
	}

	return nil
}

func scanPlugins() error {
	if verbose {
		fmt.Printf("Scanning plugins in: %s\n", scanDirectory)
	}

	results := scanner.ScanPlugins(scanDirectory)

	if len(results) == 0 {
		fmt.Println("✅ No risky plugins found")
		return nil
	}

	fmt.Printf("\nFound %d plugin(s) with potential risks:\n\n", len(results))

	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.PluginName)
		fmt.Printf("   Path: %s\n", result.Path)
		fmt.Printf("   Risk: %s\n", result.Risk)
		fmt.Printf("   Reason: %s\n", result.Reason)
		if len(result.Detected) > 0 {
			fmt.Printf("   Detected: %v\n", result.Detected)
		}
		fmt.Println()
	}

	// Print summary
	summary := scanner.GetPluginRiskSummary(results)
	fmt.Println("Risk Summary:")
	fmt.Printf("  High: %d\n", summary["HIGH"])
	fmt.Printf("  Medium: %d\n", summary["MEDIUM"])
	fmt.Printf("  Low: %d\n", summary["LOW"])

	return nil
}
