package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is set at build time using ldflags
var version = "dev"


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the version number of AgentGuard.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("AgentGuard %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
