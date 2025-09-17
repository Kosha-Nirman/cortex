package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Kosha-Nirman/cortex/src/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cortex [domain]",
	Short: "A production-grade subdomain resolver and reconnaissance tool",
	Long: `Cortex is a comprehensive subdomain discovery tool that combines multiple
reconnaissance techniques to find subdomains for a given domain.

Features:
- DNS enumeration with custom wordlists
- Certificate Transparency log searches  
- Passive reconnaissance from multiple sources
- Brute force subdomain discovery
- Detailed markdown reports
- Cross-platform support

Example:
  cortex example.com
  cortex --threads 200 --timeout 10s example.com
  cortex --no-brute --no-passive example.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.PrintBanner()
		// Create context with cancellation
		_, cancel := context.WithCancel(context.Background())
		defer cancel()

		domain := args[0]
		fmt.Printf("Detecting subdomains for: %s\n", domain)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
