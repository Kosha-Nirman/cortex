package resolver

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Kosha-Nirman/cortex/src/domain"
)

type BruteResolver struct {
	wordlist    []string
	dnsResolver *DNSResolver
	timeout     time.Duration
}

func getDefaultWordlist() []string {
	return []string{
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
		"2019", "2020", "2021", "2022", "2023", "2024", "old", "new",
		"v1", "v2", "v3", "version", "latest", "current", "next", "prev",
		"temp", "tmp", "backup", "bak", "old", "archive", "deleted",
		"public", "private", "internal", "external", "intranet", "extranet",
		"live", "dead", "active", "inactive", "online", "offline",
		"en", "us", "uk", "ca", "au", "de", "fr", "es", "it", "jp", "cn",
		"www2", "www3", "web", "web1", "web2", "website", "site",
		"proxy", "gateway", "router", "switch", "firewall", "lb", "balancer",
		"cluster", "node", "worker", "master", "slave", "primary", "secondary",
		"east", "west", "north", "south", "central", "local", "remote",
		"dev1", "dev2", "test1", "test2", "stage1", "stage2", "prod1", "prod2",
		"jenkins", "ci", "cd", "build", "deploy", "release", "publish",
		"docker", "k8s", "kubernetes", "container", "vm", "cloud", "aws", "azure",
		"service", "microservice", "lambda", "function", "webhook", "callback",
		"queue", "worker", "job", "task", "cron", "scheduler", "timer",
		"logger", "log", "audit", "trace", "debug", "error", "warning", "info",
	}
}

func loadWordlistFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open wordlist file: %w", err)
	}
	defer file.Close()

	var wordlist []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			wordlist = append(wordlist, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading wordlist file: %w", err)
	}

	return wordlist, nil
}

func NewBruteResolver(wordlistPath string, dnsResolver *DNSResolver, timeout time.Duration) (*BruteResolver, error) {
	var wordlist []string

	if wordlistPath != "" {
		loadedWordlist, err := loadWordlistFromFile(wordlistPath)
		if err != nil {
			return nil, err
		}
		wordlist = loadedWordlist
	} else {
		wordlist = getDefaultWordlist()
	}

	return &BruteResolver{
		wordlist:    wordlist,
		dnsResolver: dnsResolver,
		timeout:     timeout,
	}, nil
}

func (b *BruteResolver) DiscoverSubdomains(ctx context.Context, domain string, threads int) ([]*domain.Subdomain, error) {
	if threads <= 0 {
		threads = 100
	}

	return b.dnsResolver.ResolveSubdomains(ctx, domain, b.wordlist, threads)
}
