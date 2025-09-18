package domain

import "time"

type ResolverConfig struct {
	Timeout       time.Duration
	Threads       int
	DNSServers    []string
	WordlistPath  string
	EnableBrute   bool
	EnableCRT     bool
	EnableDNS     bool
	EnablePassive bool
	HTTPTimeout   time.Duration
}
