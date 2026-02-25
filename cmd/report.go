package cmd

import (
	"github.com/imdlan/AIAgentGuard/internal/report"
	"github.com/imdlan/AIAgentGuard/internal/risk"
	"github.com/imdlan/AIAgentGuard/internal/scanner"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate and display security reports",
	Long: `Generate detailed security reports based on scan results.
Reports can be output in various formats for different use cases.`,
	Example: `  agent-guard report
  agent-guard report --json
  agent-guard report --compact`,
	RunE: generateReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func generateReport(cmd *cobra.Command, args []string) error {
	// Run scans
	result := scanner.RunAllScans()
	analysis := risk.Analyze(result)

	// Output report
	if jsonOutput {
		return report.PrintJSON(analysis)
	}

	// Check for compact mode flag
	compact, _ := cmd.Flags().GetBool("compact")
	if compact {
		report.PrintCompact(analysis)
		return nil
	}

	// Default console output
	report.PrintConsole(analysis)
	return nil
}

func init() {
	reportCmd.Flags().Bool("compact", false, "Show compact one-line report")
}
