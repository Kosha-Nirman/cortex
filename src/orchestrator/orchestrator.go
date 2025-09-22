package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Kosha-Nirman/cortex/src/config"
	"github.com/Kosha-Nirman/cortex/src/domain"
	"github.com/Kosha-Nirman/cortex/src/resolver"
)

type Orchestrator struct {
	config              *config.ResolverConfig
	dnsResolver         *resolver.DNSResolver
	certificateResolver *resolver.CertificateResolver
	bruteResolver       *resolver.BruteResolver
	passiveResolver     *resolver.PassiveResolver
}

func NewOrchestrator(config *config.ResolverConfig) (*Orchestrator, error) {
	dnsResolver := resolver.NewDNSResolver(config.DNSServers, config.Timeout)
	certificateResolver := resolver.NewCertificateResolver(config.HTTPTimeout)
	passiveResolver := resolver.NewPassiveResolver(config.Timeout)

	var bruteResolver *resolver.BruteResolver
	var err error
	if config.EnableBrute {
		bruteResolver, err = resolver.NewBruteResolver(config.WordlistPath, dnsResolver, config.Timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to create brute force resolver: %w", err)
		}
	}

	return &Orchestrator{
		config:              config,
		dnsResolver:         dnsResolver,
		certificateResolver: certificateResolver,
		bruteResolver:       bruteResolver,
		passiveResolver:     passiveResolver,
	}, nil

}

func (o *Orchestrator) GetDomainInfo(domain string) (map[string]any, error) {
	info := make(map[string]any)

	mxRecords, err := o.dnsResolver.GetMXRecords(domain)
	if err == nil {
		info["mx_records"] = mxRecords
	}

	nsRecords, err := o.dnsResolver.GetNSRecords(domain)
	if err == nil {
		info["ns_records"] = nsRecords
	}

	txtRecords, err := o.dnsResolver.GetTXTRecords(domain)
	if err == nil {
		info["txt_records"] = txtRecords
	}

	certInfo, err := o.certificateResolver.GetCertificateInfo(domain)
	if err == nil {
		info["certificate"] = certInfo
	}

	return info, nil
}

func (o *Orchestrator) DiscoverSubdomains(ctx context.Context, targetDomain string) (*domain.ScanResult, error) {
	startTime := time.Now()

	result := &domain.ScanResult{
		Domain:    targetDomain,
		StartTime: startTime,
	}

	subdomains := make(map[string]*domain.Subdomain)
	var mu sync.Mutex
	var wg sync.WaitGroup

	if o.config.EnableDNS {
		wg.Go(func() {
			if subs, err := o.dnsResolver.DiscoverSubdomains(ctx, targetDomain); err == nil {
				mu.Lock()
				for _, sub := range subs {
					subdomains[sub.Name] = sub
				}
				mu.Unlock()
			}
		})
	}

	if o.config.EnableCRT {
		wg.Go(func() {
			if subs, err := o.certificateResolver.DiscoverSubdomains(ctx, targetDomain); err == nil {
				mu.Lock()
				for _, sub := range subs {
					subdomains[sub.Name] = sub
				}
				mu.Unlock()
			}
		})
	}

	if o.config.EnablePassive {
		wg.Go(func() {
			if subs, err := o.passiveResolver.DiscoverSubdomains(ctx, targetDomain); err == nil {
				mu.Lock()
				for _, sub := range subs {
					subdomains[sub.Name] = sub
				}
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if o.config.EnableBrute && o.bruteResolver != nil {
		bruteSubs, err := o.bruteResolver.DiscoverSubdomains(ctx, targetDomain, o.config.Threads)
		if err == nil {
			for _, sub := range bruteSubs {
				if existing, exists := subdomains[sub.Name]; exists {
					existing.Source += ", Brute Force"
				} else {
					subdomains[sub.Name] = sub
				}
			}
		}
	}

	finalSubs := make([]*domain.Subdomain, 0, len(subdomains))
	resolverWg := sync.WaitGroup{}
	semaphore := make(chan struct{}, o.config.Threads)

	for _, sub := range subdomains {
		finalSubs = append(finalSubs, sub)

		resolverWg.Add(1)
		go func(s *domain.Subdomain) {
			defer resolverWg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := s.Resolve(); err != nil {
				fmt.Printf("failed to resolve subdomain %s: %v\n", s.Name, err)
			}
		}(sub)
	}

	resolverWg.Wait()

	result.Subdomains = finalSubs
	result.EndTime = time.Now()
	result.CalculateStats()

	return result, nil
}
