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

type CertificateResolver struct {
	client  *http.Client
	timeout time.Duration
}

type CrtShResponse struct {
	IssuerCAID        int    `json:"issuer_ca_id"`
	IssuerName        string `json:"issuer_name"`
	NameValue         string `json:"name_value"`
	MinCertId         int64  `json:"min_cert_id"`
	MinEntryTimestamp string `json:"min_entry_timestamp"`
	NotBefore         string `json:"not_before"`
	NotAfter          string `json:"not_after"`
}

func NewCertificateResolver(timeout time.Duration) *CertificateResolver {
	return &CertificateResolver{
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

func (c *CertificateResolver) searchCrtSh(ctx context.Context, targetDomain string) ([]*domain.Subdomain, error) {
	encodedDomain := url.QueryEscape("%." + targetDomain)
	apiURL := fmt.Sprintf("https://crt.sh/?q=%s&output=json", encodedDomain)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Cortex-SubdomainResolver/1.0.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("crt.sh API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var crtShResponse []CrtShResponse
	if err := json.Unmarshal(body, &crtShResponse); err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var results []*domain.Subdomain

	for _, cert := range crtShResponse {
		names := strings.SplitSeq(cert.NameValue, "\n")
		for name := range names {
			name = strings.ToLower(strings.TrimSpace(name))

			if name == "" || seen[name] {
				continue
			}

			if strings.Contains(name, "*") {
				continue
			}

			if strings.HasSuffix(name, "."+targetDomain) || name == targetDomain {
				seen[name] = true
				sub := domain.NewSubdomain(name, "Certificate Transparency")
				results = append(results, sub)
			}
		}
	}

	return results, nil
}

func (c *CertificateResolver) DiscoverSubdomains(ctx context.Context, targetDomain string) ([]*domain.Subdomain, error) {
	subdomains := make(map[string]*domain.Subdomain)

	crtSubdomains, err := c.searchCrtSh(ctx, targetDomain)
	if err == nil {
		for _, sub := range crtSubdomains {
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

func (c *CertificateResolver) GetCertificateInfo(domain string) (map[string]any, error) {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]

	info := map[string]any{
		"subject":    cert.Subject.String(),
		"issuer":     cert.Issuer.String(),
		"not_before": cert.NotBefore,
		"not_after":  cert.NotAfter,
		"dns_names":  cert.DNSNames,
		"serial":     cert.SerialNumber.String(),
		"is_ca":      cert.IsCA,
	}

	return info, nil
}
