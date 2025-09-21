package orchestrator

import (
	"fmt"

	"github.com/Kosha-Nirman/cortex/src/config"
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
