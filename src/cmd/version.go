package cmd

import (
	"github.com/Kosha-Nirman/cortex/src/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display version, build date, and git commit information for Cortex",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintBanner()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
