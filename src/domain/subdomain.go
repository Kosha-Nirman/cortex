package domain

import (
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
