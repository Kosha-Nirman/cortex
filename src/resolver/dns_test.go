package resolver

import (
	"context"
	"testing"
	"time"
)

func TestNewDNSResolver(t *testing.T) {
	resolver := NewDNSResolver(nil, 5*time.Second)

	if len(resolver.servers) == 0 {
		t.Error("DNS resolver should have default servers when none provided")
	}

	if resolver.timeout != 5*time.Second {
		t.Error("DNS resolver timeout not set correctly")
	}

	customServers := []string{"1.1.1.1:53", "9.9.9.9:53"}
	resolver2 := NewDNSResolver(customServers, 3*time.Second)

	if len(resolver2.servers) != 2 {
		t.Error("DNS resolver should use custom servers when provided")
	}

	if resolver2.servers[0] != customServers[0] {
		t.Error("DNS resolver should preserve custom server order")
	}
}

func TestDNSResolverDiscoverSubdomains(t *testing.T) {
	resolver := NewDNSResolver(nil, 5*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subdomains, err := resolver.DiscoverSubdomains(ctx, "google.com")
	if err != nil {
		t.Fatalf("DiscoverSubdomains failed: %v", err)
	}

	if len(subdomains) == 0 {
		t.Log("No subdomains found - this may be expected for some domains")
	}

	for _, sub := range subdomains {
		if !sub.IsValid() {
			t.Errorf("Invalid subdomain found: %s", sub.Name)
		}

		if sub.Source != "DNS" {
			t.Errorf("Expected source to be 'DNS', got '%s'", sub.Source)
		}
	}
}
