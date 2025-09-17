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
	Short: "Subdomain detector CLI tool",
	Long:  `Cortex is a command-line tool to detect subdomains associated with a domain name.`,
	Args:  cobra.ExactArgs(1),
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
