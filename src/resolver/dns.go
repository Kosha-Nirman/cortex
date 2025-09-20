package resolver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Kosha-Nirman/cortex/src/domain"
)

type DNSResolver struct {
	servers []string
	timeout time.Duration
}

func NewDNSResolver(servers []string, timeout time.Duration) *DNSResolver {
	if len(servers) == 0 {
		servers = []string{"8.8.8.8:53", "1.1.1.1:53", "9.9.9.9:53"}
	}

	return &DNSResolver{
		servers: servers,
		timeout: timeout,
	}
}

func (d *DNSResolver) ResolveSubdomains(ctx context.Context, rootDomain string, subdomains []string, threads int) ([]*domain.Subdomain, error) {
	if threads <= 0 {
		threads = 50
	}

	semaphore := make(chan struct{}, threads)
	results := make([]*domain.Subdomain, 0, len(subdomains))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, subdomain := range subdomains {
		wg.Add(1)

		go func(sub string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				return
			default:
			}

			fullDomain := fmt.Sprintf("%s.%s", sub, rootDomain)
			resolved := domain.NewSubdomain(fullDomain, "DNS")

			if err := resolved.Resolve(); err == nil && resolved.IsActive {
				mu.Lock()
				results = append(results, resolved)
				mu.Unlock()
			}
		}(subdomain)
	}

	wg.Wait()
	return results, nil
}

func (d *DNSResolver) DiscoverSubdomains(ctx context.Context, rootDomain string) ([]*domain.Subdomain, error) {
	var results []*domain.Subdomain

	commonSubs := []string{
		"www", "mail", "ftp", "admin", "api", "blog", "dev", "test", "staging",
		"app", "cdn", "assets", "static", "media", "images", "js", "css",
		"support", "help", "docs", "forum", "shop", "store", "news", "portal",
		"vpn", "secure", "ssl", "mx", "mx1", "mx2", "smtp", "pop", "imap",
		"ns", "ns1", "ns2", "dns", "dns1", "dns2", "server", "host", "git",
		"gitlab", "github", "svn", "repo", "download", "files", "upload",
		"beta", "alpha", "demo", "preview", "staging", "prod", "production",
		"db", "database", "mysql", "postgres", "redis", "cache", "backup",
		"monitor", "status", "health", "metrics", "analytics", "stats",
		"auth", "oauth", "sso", "login", "account", "profile", "dashboard",
		"mobile", "m", "wap", "touch", "tablet", "ios", "android",
		"email", "webmail", "mail2", "mailserver", "exchange", "outlook",
		"crm", "erp", "hr", "finance", "accounting", "billing", "invoice",
		"chat", "im", "messenger", "video", "voice", "webrtc", "conference",
		"search", "find", "directory", "catalog", "archive", "library",
		"wiki", "kb", "knowledgebase", "faq", "terms", "privacy", "legal",
	}

	discovered, err := d.ResolveSubdomains(ctx, rootDomain, commonSubs, 50)
	if err != nil {
		return nil, err
	}

	results = append(results, discovered...)

	return results, nil
}

func (d *DNSResolver) GetMXRecords(domain string) ([]string, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, mx := range mxRecords {
		host := strings.TrimSuffix(mx.Host, ".")
		results = append(results, host)
	}

	return results, nil
}

func (d *DNSResolver) GetNSRecords(domain string) ([]string, error) {
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, ns := range nsRecords {
		host := strings.TrimSuffix(ns.Host, ".")
		results = append(results, host)
	}

	return results, nil
}

func (d *DNSResolver) GetTXTRecords(domain string) ([]string, error) {
	return net.LookupTXT(domain)
}
