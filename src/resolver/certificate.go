package resolver

import (
	"crypto/tls"
	"net/http"
	"time"
)

type CertificateResovler struct {
	client  *http.Client
	timeout time.Duration
}

func NewCertificateResolver(timeout time.Duration) *CertificateResovler {
	return &CertificateResovler{
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
