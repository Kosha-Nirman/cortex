package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rivetron/cortex/src/config"
	"github.com/rivetron/cortex/src/orchestrator"
	"github.com/rivetron/cortex/src/report"
	"github.com/rivetron/cortex/src/utils"
	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	threads       int
	timeout       time.Duration
	httpTimeout   time.Duration
	wordlistPath  string
	outputDir     string
	noColor       bool
	enableCRT     bool
	enableDNS     bool
	enableBrute   bool
	enablePassive bool
	dnsServers    []string
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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		targetDomain := args[0]

		resolverConfig := &config.ResolverConfig{
			Timeout:       timeout,
			Threads:       threads,
			DNSServers:    dnsServers,
			WordlistPath:  wordlistPath,
			EnableBrute:   enableBrute && !cmd.Flags().Changed("no-brute"),
			EnableCRT:     enableCRT && !cmd.Flags().Changed("no-crt"),
			EnableDNS:     enableDNS && !cmd.Flags().Changed("no-dns"),
			EnablePassive: enablePassive && !cmd.Flags().Changed("no-passive"),
			HTTPTimeout:   httpTimeout,
		}

		fmt.Printf("🎯 Target: %s\n", color.CyanString(targetDomain))
		fmt.Println()
		fmt.Printf("⚙️  Configuration:\n")
		fmt.Printf("   Threads: %d\n", threads)
		fmt.Printf("   Timeout: %s\n", timeout)
		fmt.Printf("   DNS: %v | CRT: %v | Brute: %v | Passive: %v\n",
			resolverConfig.EnableDNS, resolverConfig.EnableCRT,
			resolverConfig.EnableBrute, resolverConfig.EnablePassive)
		fmt.Println()

		fmt.Printf("🔍 Starting subdomain discovery...\n\n")

		orchestrator, err := orchestrator.NewOrchestrator(resolverConfig)
		if err != nil {
			return fmt.Errorf("failed to create resolver: %w", err)
		}

		result, err := orchestrator.DiscoverSubdomains(ctx, targetDomain)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		utils.PrintResults(result)

		reportGen, err := report.NewMarkdownGenerator(outputDir)
		if err != nil {
			return fmt.Errorf("failed to create report generator: %w", err)
		}

		reportPath, err := reportGen.GenerateReport(result)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}

		fmt.Printf("\n📋 Report saved: %s\n", color.GreenString(reportPath))

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

func initConfig() {
	if noColor {
		color.NoColor = true
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.cortex.yaml)")
	rootCmd.Flags().IntVarP(&threads, "threads", "t", 100, "Number of concurrent threads")
	rootCmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "DNS resolution timeout")
	rootCmd.Flags().DurationVar(&httpTimeout, "http-timeout", 10*time.Second, "HTTP request timeout")
	rootCmd.Flags().StringVarP(&wordlistPath, "wordlist", "w", "", "Path to custom wordlist file")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: ~/Downloads/cortex-reports)")
	rootCmd.Flags().StringSliceVar(&dnsServers, "dns-servers", []string{"8.8.8.8:53", "1.1.1.1:53"}, "DNS servers to use (comma-separated)")

	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	rootCmd.Flags().BoolVar(&enableCRT, "crt", true, "Enable Certificate Transparency log searches")
	rootCmd.Flags().BoolVar(&enableDNS, "dns", true, "Enable DNS enumeration")
	rootCmd.Flags().BoolVar(&enableBrute, "brute", true, "Enable brute force subdomain discovery")
	rootCmd.Flags().BoolVar(&enablePassive, "passive", true, "Enable passive reconnaissance")

	rootCmd.Flags().BoolVar(&enableCRT, "no-crt", true, "Disable Certificate Transparency (opposite of --crt)")
	rootCmd.Flags().BoolVar(&enableDNS, "no-dns", true, "Disable DNS enumeration (opposite of --dns)")
	rootCmd.Flags().BoolVar(&enableBrute, "no-brute", true, "Disable brute force (opposite of --brute)")
	rootCmd.Flags().BoolVar(&enablePassive, "no-passive", true, "Disable passive recon (opposite of --passive)")

	rootCmd.MarkFlagsMutuallyExclusive("crt", "no-crt")
	rootCmd.MarkFlagsMutuallyExclusive("dns", "no-dns")
	rootCmd.MarkFlagsMutuallyExclusive("brute", "no-brute")
	rootCmd.MarkFlagsMutuallyExclusive("passive", "no-passive")
}
