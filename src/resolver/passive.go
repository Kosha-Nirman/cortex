package resolver

import (
	"crypto/tls"
	"net/http"
	"time"
)

type PassiveResolver struct {
	client  *http.Client
	timeout time.Duration
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
