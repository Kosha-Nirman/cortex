package helper

import (
	"fmt"
	"strings"
)

var commonSubdomains = []string{
	"www", "mail", "ftp", "blog", "shop", "api", "dev", "staging", "test",
	"admin", "login", "secure", "portal", "client", "support", "feedback",
}

// FindSubdomains checks if common subdomains exist for a given domain
func FindSubdomains(domain string) []string {
	var foundDomains []string
	for _, subdomain := range commonSubdomains {
		fullDomain := strings.ToLower(subdomain + "." + domain)
		if isAvailable(fullDomain) {
			foundDomains = append(foundDomains, fullDomain)
		}
	}
	return foundDomains
}

// TODO: Implement actual DNS lookup to verify subdomain existence
func isAvailable(domain string) bool {
	fmt.Printf("Checking availability of: %s\n", domain)
	return false // Placeholder - always returns false
}
