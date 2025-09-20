package resolver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Kosha-Nirman/cortex/src/domain"
)

type HackerTargetResponse struct {
	Subdomains []string `json:"subdomains"`
}

type PassiveResolver struct {
	client  *http.Client
	timeout time.Duration
}

type VirusTotalResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			LastAnalysisDate int64 `json:"last_analysis_date"`
		} `json:"attributes"`
	} `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

func NewPassiveResolver(timeout time.Duration) *PassiveResolver {
	return &PassiveResolver{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		timeout: timeout,
	}
}

func (p *PassiveResolver) searchHackerTarget(ctx context.Context, targetDomain string) ([]*domain.Subdomain, error) {
	apiURL := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", url.QueryEscape(targetDomain))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Cortex-SubdomainResolver/1.0.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hackertarget API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*domain.Subdomain
	lines := strings.SplitSeq(string(body), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "error") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) >= 1 {
			subdomain := strings.ToLower(strings.TrimSpace(parts[0]))
			if subdomain != "" && strings.HasSuffix(subdomain, "."+targetDomain) {
				sub := domain.NewSubdomain(subdomain, "HackerTarget")
				results = append(results, sub)
			}
		}
	}

	return results, nil
}

func (p *PassiveResolver) searchThreatCrowd(ctx context.Context, targetDomain string) ([]*domain.Subdomain, error) {
	apiURL := fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", url.QueryEscape(targetDomain))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Cortex-SubdomainResolver/1.0.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("threatcrowd API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Subdomains   []string `json:"subdomains"`
		ResponseCode string   `json:"response_code"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.ResponseCode != "1" {
		return nil, fmt.Errorf("threatcrowd API returned error code: %s", response.ResponseCode)
	}

	var results []*domain.Subdomain
	for _, subdomain := range response.Subdomains {
		subdomain = strings.ToLower(strings.TrimSpace(subdomain))
		if subdomain != "" {
			sub := domain.NewSubdomain(subdomain, "ThreatCrowd")
			results = append(results, sub)
		}
	}

	return results, nil
}

func (p *PassiveResolver) DiscoverSubDomains(ctx context.Context, targetDomain string) ([]*domain.Subdomain, error) {
	subdomains := make(map[string]*domain.Subdomain)

	hackertargetSubs, err := p.searchHackerTarget(ctx, targetDomain)
	if err == nil {
		for _, sub := range hackertargetSubs {
			if sub.IsValid() {
				subdomains[sub.Name] = sub
			}
		}
	}

	threatcrowdSubs, err := p.searchThreatCrowd(ctx, targetDomain)
	if err == nil {
		for _, sub := range threatcrowdSubs {
			if sub.IsValid() {
				subdomains[sub.Name] = sub
			}
		}
	}

	results := make([]*domain.Subdomain, 0, len(subdomains))
	for _, sub := range subdomains {
		results = append(results, sub)
	}

	return results, nil
}
