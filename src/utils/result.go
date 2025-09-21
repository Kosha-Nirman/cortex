package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Kosha-Nirman/cortex/src/domain"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func PrintResults(result *domain.ScanResult) {
	fmt.Printf("✅ Scan completed in %s\n", color.GreenString(result.Duration))
	fmt.Printf("📊 Results:\n")
	fmt.Printf("   Total subdomains found: %s\n", color.CyanString(fmt.Sprintf("%d", result.TotalFound)))
	fmt.Printf("   Active subdomains: %s\n", color.GreenString(fmt.Sprintf("%d", result.ActiveSubs)))
	fmt.Printf("   Sources used: %s\n", color.YellowString(fmt.Sprintf("%s", result.Sources)))

	if result.ActiveSubs > 0 {
		color.Magenta("\n🎯 Active Subdomains:")

		table := tablewriter.NewWriter(os.Stdout)

		// Add Header
		table.Header([]string{"Subdomain", "IP Address", "Source"})

		// Build rows
		for _, sub := range result.Subdomains {
			if sub.IsActive {
				ip := ""
				if len(sub.IPAddresses) > 0 {
					ip = sub.IPAddresses[0]
				}

				// Colorize source conditionally
				src := sub.Source
				if strings.Contains(strings.ToLower(src), "hackertarget") {
					src = color.RedString(src) // make HackerTarget red
				} else {
					src = color.YellowString(src)
				}

				row := []string{
					color.GreenString(sub.Name),
					color.CyanString(ip),
					src,
				}
				ok := table.Append(row)
				if ok != nil {
					fmt.Printf("⚠️ Failed to append row: %+v\n", row)
				}

			}
		}

		// Add footer with total
		table.Footer([]string{"", "Active", strconv.Itoa(result.ActiveSubs)})

		// Render table
		if err := table.Render(); err != nil {
			fmt.Printf("⚠️ Failed to render table: %v\n", err)
		}
	}
}
