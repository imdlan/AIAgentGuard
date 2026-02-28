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

		// Group details by category for better readability
		byCategory := groupDetailsByCategory(report.Details)

		for _, category := range []string{"filesystem", "shell", "network", "secrets"} {
			if details, ok := byCategory[category]; ok {
				printDetailedCategory(category, details)
			}
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


// groupDetailsByCategory groups risk details by their category
func groupDetailsByCategory(details []model.RiskDetail) map[string][]model.RiskDetail {
	byCategory := make(map[string][]model.RiskDetail)
	for _, detail := range details {
		byCategory[detail.Category] = append(byCategory[detail.Category], detail)
	}
	return byCategory
}

// printDetailedCategory prints detailed information for a specific category
func printDetailedCategory(category string, details []model.RiskDetail) {
	fmt.Println()
	switch category {
	case "filesystem":
		fmt.Printf("📁 Filesystem Risk\n")
		for _, detail := range details {
			fmt.Printf("  %s %s\n", getRiskSymbol(detail.Type), detail.Description)
			if len(detail.Details.AffectedPaths) > 0 {
				for _, path := range detail.Details.AffectedPaths {
					fmt.Printf("     └─ %s\n", path.Path)
					fmt.Printf("        └─ Risk: %s\n", path.RiskReason)
					fmt.Printf("        └─ Permission: %s\n", path.Permission)
					fmt.Printf("        └─ Writable: %v\n", path.IsWritable)
				}
			}
			// Print remediation if available
			if detail.Remediation.Summary != "" {
				printRemediation(detail.Remediation)
			}
		}

	case "shell":
		fmt.Printf("💻 Shell Risk\n")
		for _, detail := range details {
			fmt.Printf("  %s %s\n", getRiskSymbol(detail.Type), detail.Description)
			if detail.Details.ShellAvailable != "" {
				fmt.Printf("     └─ Available Shells: %s\n", detail.Details.ShellAvailable)
			}
			if detail.Details.HasSudoAccess {
				fmt.Printf("     └─ ⚠️  Sudo Access: ENABLED\n")
				if detail.Details.SudoSource != "" {
					fmt.Printf("        └─ Source: %s\n", detail.Details.SudoSource)
				}
				if len(detail.Details.SudoRules) > 0 {
					fmt.Printf("        └─ Rules:\n")
					for _, rule := range detail.Details.SudoRules {
						fmt.Printf("           - %s\n", rule)
					}
				}
			}
			// Print remediation if available
			if detail.Remediation.Summary != "" {
				printRemediation(detail.Remediation)
			}
		}

	case "network":
		fmt.Printf("🌐 Network Risk\n")
		for _, detail := range details {
			fmt.Printf("  %s %s\n", getRiskSymbol(detail.Type), detail.Description)
			if len(detail.Details.OpenPorts) > 0 {
				for _, port := range detail.Details.OpenPorts {
					fmt.Printf("     └─ Port %d/%s\n", port.Port, port.Protocol)
					fmt.Printf("        └─ Service: %s\n", port.Service)
					fmt.Printf("        └─ Risk: %s\n", port.RiskReason)
				}
			}
		}

	case "secrets":
		fmt.Printf("🔑 Secrets Risk\n")
		for _, detail := range details {
			fmt.Printf("  %s %s\n", getRiskSymbol(detail.Type), detail.Description)
			if detail.Path != "" {
				fmt.Printf("     └─ Exposed: %s\n", detail.Path)
			}
			if len(detail.Details.ExposedSecrets) > 0 {
				for _, secret := range detail.Details.ExposedSecrets {
					fmt.Printf("     └─ %s: %s\n", secret.Type, secret.Value)
					fmt.Printf("        └─ Location: %s\n", secret.Location)
				}
			}
		}
	}
}

// printRemediation prints remediation steps
func printRemediation(remediation model.RemediationInfo) {
	fmt.Println()
	fmt.Printf("  💡 Remediation: %s\n", remediation.Summary)
	if len(remediation.Steps) > 0 {
		fmt.Println("  Steps:")
		for _, step := range remediation.Steps {
			fmt.Printf("     %d. %s\n", step.Step, step.Action)
			fmt.Printf("        Command: %s\n", step.Command)
			fmt.Printf("        Explanation: %s\n", step.Explanation)
		}
	}
	if len(remediation.Commands) > 0 {
		fmt.Println("  Commands to run:")
		for _, cmd := range remediation.Commands {
			fmt.Printf("     $ %s\n", cmd)
		}
	}
	fmt.Printf("  Priority: %s\n", remediation.Priority)
}
