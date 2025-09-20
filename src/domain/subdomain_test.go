package domain

import (
	"testing"
)

func TestNewSubdomain(t *testing.T) {
	name := "test.example.com"
	source := "Test"

	sub := NewSubdomain(name, source)

	if sub.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, sub.Name)
	}

	if sub.Source != source {
		t.Errorf("Expected Source to be %s, got %s", source, sub.Source)
	}

	if sub.IsActive {
		t.Error("New subdomain should not be active by default")
	}

	if len(sub.IPAddresses) != 0 {
		t.Error("Expected IPAddresses to be empty")
	}

	if len(sub.CNAMEs) != 0 {
		t.Error("Expected CNAMEs to be empty")
	}
}

func TestSubdomainIsValid(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{"valid domain", "test.example.com", true},
		{"valid single level", "example.com", true},
		{"empty domain", "", false},
		{"single part", "example", false},
		{"too long part", "this-is-a-very-long-subdomain-name-that-exceeds-the-maximum-allowed-length-for-a-single-part-of-a-domain-name-which-is-63-characters.example.com", false},
		{"too long overall", "verylongsubdomainnamethatshouldexceedthemaximumlengthallowedforadomainname.verylongsubdomainnamethatshouldexceedthemaximumlengthallowedforadomainname.verylongsubdomainnamethatshouldexceedthemaximumlengthallowedforadomainname.example.com", false},
		{"empty part", "test..example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := NewSubdomain(tt.domain, "Test")

			if sub.IsValid() != tt.expected {
				t.Errorf("Expected IsValid to return %v for domain %s", tt.expected, tt.domain)
			}
		})
	}
}
