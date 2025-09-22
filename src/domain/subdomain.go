package domain

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Subdomain struct {
	Name        string    `json:"name"`
	IPAddresses []string  `json:"ip_addresses"`
	CNAMEs      []string  `json:"cnames"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp"`
	IsActive    bool      `json:"is_active"`
	HTTPStatus  int       `json:"http_status"`
	Title       string    `json:"title"`
}

type ScanResult struct {
	Domain     string       `json:"domain"`
	Subdomains []*Subdomain `json:"subdomains"`
	StartTime  time.Time    `json:"start_time"`
	EndTime    time.Time    `json:"end_time"`
	Duration   string       `json:"duration"`
	TotalFound int          `json:"total_found"`
	ActiveSubs int          `json:"active_subs"`
	Sources    []string     `json:"sources"`
}

func NewSubdomain(name, source string) *Subdomain {
	return &Subdomain{
		Name:        strings.ToLower(strings.TrimSpace(name)),
		Source:      source,
		Timestamp:   time.Now(),
		IPAddresses: []string{},
		CNAMEs:      []string{},
	}
}

func (s *Subdomain) IsValid() bool {
	if s.Name == "" {
		return false
	}

	parts := strings.Split(s.Name, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if part == "" || len(part) > 63 {
			return false
		}
	}

	return len(s.Name) <= 253
}

func (s *Subdomain) Resolve() error {
	if s.Name == "" {
		return fmt.Errorf("subdomain name cannot be empty")
	}

	ips, err := net.LookupIP(s.Name)
	if err != nil {
		s.IsActive = false
		return err
	}

	s.IPAddresses = make([]string, 0, len(ips))
	for _, ip := range ips {
		s.IPAddresses = append(s.IPAddresses, ip.String())
	}

	cname, err := net.LookupCNAME(s.Name)
	if err == nil && cname != s.Name+"." {
		s.CNAMEs = append(s.CNAMEs, strings.TrimSuffix(cname, "."))
	}

	s.IsActive = len(s.IPAddresses) > 0
	return nil
}

func (sr *ScanResult) CalculateStats() {
	sr.TotalFound = len(sr.Subdomains)
	sr.ActiveSubs = 0

	sourcesMap := make(map[string]bool)
	for _, sub := range sr.Subdomains {
		if sub.IsActive {
			sr.ActiveSubs++
		}
		sourcesMap[sub.Source] = true
	}

	sr.Sources = make([]string, 0, len(sourcesMap))
	for source := range sourcesMap {
		sr.Sources = append(sr.Sources, source)
	}

	sr.Duration = sr.EndTime.Sub(sr.StartTime).String()
}
