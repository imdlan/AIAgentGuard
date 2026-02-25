package cmd

import (
	"fmt"
	"os"

	"github.com/imdlan/AIAgentGuard/internal/policy"
	"github.com/imdlan/AIAgentGuard/internal/sandbox"
	"github.com/imdlan/AIAgentGuard/internal/security"
	"github.com/spf13/cobra"
)

var (
	runStrict      bool
	runNoNetwork   bool
	runPromptCheck bool
)

var runCmd = &cobra.Command{
	Use:   "run <command>",
	Short: "Run a command in a sandboxed environment",
	Long: `Execute a command in a sandboxed environment with restricted permissions.
This helps prevent AI agents from executing dangerous operations.`,
	Example: `  agent-guard run "ls -la"
  agent-guard run --strict "curl https://api.example.com"
  agent-guard run --check-prompt "echo test"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSandboxed,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVarP(&runStrict, "strict", "s", false, "Use strict sandbox mode (clear env, no network)")
	runCmd.Flags().BoolVarP(&runNoNetwork, "no-network", "n", false, "Disable network access")
	runCmd.Flags().BoolVarP(&runPromptCheck, "check-prompt", "p", false, "Check for prompt injection before running")
}

func runSandboxed(cmd *cobra.Command, args []string) error {
	command := ""

	if len(args) == 1 {
		command = args[0]
	} else {
		// Join multiple arguments with proper handling
		command = args[0]
		for i := 1; i < len(args); i++ {
			command += " " + args[i]
		}
	}

	// Check for prompt injection if requested
	if runPromptCheck {
		if isMalicious, reasons := security.CheckSecurity(command); isMalicious {
			fmt.Fprintln(os.Stderr, "🛑 Security Check Failed:")
			for _, reason := range reasons {
				fmt.Fprintf(os.Stderr, "  - %s\n", reason)
			}
			return fmt.Errorf("command blocked by security policy")
		}
	}

	// Load policy configuration
	cfg, err := policy.LoadConfig(configFile)
	if err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: Could not load policy config: %v\n", err)
	}

	// Check if command is allowed by policy
	if cfg != nil {
		if !policy.IsCommandAllowed(command, cfg) {
			return fmt.Errorf("command blocked by policy: not in allow list or explicitly denied")
		}
	}

	// Check for dangerous commands
	if err := sandbox.GuardCommand(command); err != nil {
		fmt.Fprintln(os.Stderr, "🛑 Dangerous Command Detected:")
		fmt.Fprintln(os.Stderr, err)
		return fmt.Errorf("command blocked: dangerous pattern detected")
	}

	// Configure sandbox
	sandboxCfg := sandbox.GetDefaultConfig()
	if runStrict || (cfg != nil && cfg.Sandbox.DisableNetwork) || runNoNetwork {
		sandboxCfg = sandbox.GetStrictConfig()
	}

	if verbose {
		fmt.Printf("Running command in sandbox mode:\n")
		fmt.Printf("  Command: %s\n", command)
		fmt.Printf("  Strict mode: %v\n", runStrict)
		fmt.Printf("  Network disabled: %v\n", sandboxCfg.DisableNetwork)
		fmt.Printf("  Working directory: %s\n", sandboxCfg.WorkingDir)
		fmt.Println()
	}

	// Run the command in sandbox
	if err := sandbox.RunSandboxed(command, sandboxCfg); err != nil {
		return fmt.Errorf("sandbox execution failed: %w", err)
	}

	return nil
}
