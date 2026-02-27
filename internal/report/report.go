package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// PrintConsole formats and prints the scan report to console
func PrintConsole(report model.ScanReport) {
	fmt.Println()
	printBanner()
	fmt.Println()

	// Overall risk with color indicator
	riskSymbol := getRiskSymbol(report.Overall)
	fmt.Printf("Overall Risk: %s %s\n", riskSymbol, report.Overall)
	fmt.Println()

	// Permission breakdown
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Permission Breakdown:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	printPermissionStatus("Filesystem Access", report.Results.Filesystem)
	printPermissionStatus("Shell Execution", report.Results.Shell)
	printPermissionStatus("Network Access", report.Results.Network)
	printPermissionStatus("Secrets Access", report.Results.Secrets)

	// Detailed findings
	if len(report.Details) > 0 {
		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("Detailed Findings:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for i, detail := range report.Details {
			fmt.Printf("%d. [%s] %s", i+1, detail.Type, detail.Description)
			if detail.Path != "" {
				fmt.Printf(" (%s)", detail.Path)
			}
			if detail.Category != "" {
				fmt.Printf(" [%s]", detail.Category)
			}
			fmt.Println()
		}
	}

	// Recommendations
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Recommendations:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	printRecommendations(report)

	fmt.Println()
}

// printBanner prints the ASCII art banner
func printBanner() {
	printFullBanner()
}

// printFullBanner prints the ASCII art banner
func printFullBanner() {
	fmt.Println("  █████╗ ██╗     █████╗  ██████╗ ███████╗███╗   ██╗████████╗ ██████╗ ██╗   ██╗ █████╗ ██████╗ ██████╗ ")
	fmt.Println(" ██╔══██╗██║    ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝██╔════╝ ██║   ██║██╔══██╗██╔══██╗██╔══██╗")
	fmt.Println(" ███████║██║    ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   ██║  ███╗██║   ██║███████║██████╔╝██║  ██║")
	fmt.Println(" ██╔══██╗██║    ██╔══██╗██║   ██║██╔══╝  ██║╚██╗██║   ██║   ██║   ██║██║   ██║██╔══██╗██╔══██╗██║  ██║")
	fmt.Println(" ██║  ██║██║    ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   ╚██████╔╝╚██████╔╝██║  ██║██║  ██║██████╔╝")
	fmt.Println(" ╚═╝  ╚═╝╚═╝    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝")
	fmt.Println()
	fmt.Println("                                 🛡️  Security Scan Report v1.0                        ")
}

// printPermissionStatus prints a permission status with indicator
func printPermissionStatus(label string, level model.RiskLevel) {
	symbol := getRiskSymbol(level)
	fmt.Printf("  %s %s: %s\n", symbol, label, level)
}

// getRiskSymbol returns a unicode symbol for the risk level
func getRiskSymbol(level model.RiskLevel) string {
	switch level {
	case model.Low:
		return "✅"
	case model.Medium:
		return "⚠️"
	case model.High:
		return "🔶"
	case model.Critical:
		return "🛑"
	default:
		return "❓"
	}
}

// printRecommendations prints security recommendations based on the scan results
func printRecommendations(report model.ScanReport) {
	if report.Results.Shell == model.Critical || report.Results.Shell == model.High {
		fmt.Println("  • Consider running AI agents in a sandboxed environment")
		fmt.Println("  • Use 'agent-guard run <command>' for safe execution")
	}

	if report.Results.Filesystem == model.High {
		fmt.Println("  • Restrict file access using policy configuration")
		fmt.Println("  • Create .agent-guard.yaml with deny rules")
	}

	if report.Results.Secrets == model.High {
		fmt.Println("  • Use environment variable blocking in policy config")
		fmt.Println("  • Consider using secret management tools")
	}

	if report.Results.Network == model.Medium {
		fmt.Println("  • Restrict network access in sandbox mode")
		fmt.Println("  • Use 'disable_network: true' in policy config")
	}

	if report.Overall == model.Low {
		fmt.Println("  • Your environment is relatively secure")
		fmt.Println("  • Continue monitoring for changes")
	}
}

// PrintJSON outputs the report as JSON
func PrintJSON(report model.ScanReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// PrintCompact outputs a compact one-line summary
func PrintCompact(report model.ScanReport) {
	parts := []string{
		fmt.Sprintf("Risk:%s", report.Overall),
		fmt.Sprintf("FS:%s", report.Results.Filesystem),
		fmt.Sprintf("Shell:%s", report.Results.Shell),
		fmt.Sprintf("Net:%s", report.Results.Network),
		fmt.Sprintf("Secrets:%s", report.Results.Secrets),
	}
	fmt.Println(strings.Join(parts, " | "))
}
