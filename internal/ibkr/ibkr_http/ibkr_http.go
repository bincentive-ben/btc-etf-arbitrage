package ibkr_http

import (
	"crypto/tls"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/go-resty/resty/v2"
)

type IBKRHttpClient struct {
	client *resty.Client
	url    string
}

// NewIBKRHttpClient creates a new instance of IBKRClient
func NewIBKRHttpClient() *IBKRHttpClient {

	httpClient := resty.New()
	httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return &IBKRHttpClient{
		client: httpClient,
		url:    config.ArbitrageConfig.AppConfig.IBKR.RestEndpoint,
	}
}
