package ibkr

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/go-resty/resty/v2"
)

type IBKRClient struct {
	client  *resty.Client
	baseUrl string
	isPaper bool
}

// NewIBKRClient creates a new instance of IBKRClient
func NewIBKRClient() *IBKRClient {

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return &IBKRClient{
		client:  client,
		baseUrl: config.AppConfig.IBKR.Endpoint,
		isPaper: config.AppConfig.IBKR.IsPaper,
	}
}

// GetAuthenticationStatus get authentication status to the Brokerage system.
func (c *IBKRClient) GetAuthenticationStatus() (bool, error) {
	type Response struct {
		Authenticated bool   `json:"authenticated"`
		Competing     bool   `json:"competing"`
		Connected     bool   `json:"connected"`
		Message       string `json:"message"`
		MAC           string `json:"MAC"`
		ServerInfo    struct {
			ServerName    string `json:"serverName"`
			ServerVersion string `json:"serverVersion"`
		} `json:"serverInfo"`
		HardwareInfo string `json:"hardware_info"`
		Fail         string `json:"fail"`
	}

	resp, err := c.client.R().
		Post(fmt.Sprintf("%s/iserver/auth/status", c.baseUrl))

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("failed to get authentication status")
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return false, err
	}

	return response.Authenticated, nil
}

// Logout logout from the Brokerage system.
func (c *IBKRClient) Logout() (string, error) {
	type Response struct {
		Status string `json:"status"`
	}

	resp, err := c.client.R().
		Post(fmt.Sprintf("%s/logout", c.baseUrl))

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to log out")
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}

// ValidateSession validates the current session for the SSO user.
func (c *IBKRClient) ValidateSession() (bool, error) {
	type Response struct {
		LoginType int    `json:"LOGIN_TYPE"`
		UserName  string `json:"USER_NAME"`
		UserID    int    `json:"USER_ID"`
		Expire    int    `json:"expire"`
		Result    bool   `json:"RESULT"`
		AuthTime  int    `json:"AUTH_TIME"`
	}

	resp, err := c.client.R().
		Get(fmt.Sprintf("%s/sso/validate", c.baseUrl))

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("failed to validate session")
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return false, err
	}

	return response.Result, nil
}

// GetMarketData fetches market data for a specific symbol (BTC or ETFs)
func (c *IBKRClient) GetMarketData(symbol string) (float64, error) {
	var response struct {
		Price float64 `json:"price"`
	}

	resp, err := c.client.R().
		SetResult(&response).
		Get(fmt.Sprintf("%s/marketdata/%s", c.baseUrl, symbol))

	if err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, fmt.Errorf("failed to fetch market data for %s", symbol)
	}

	return response.Price, nil
}

// Tickle sends a tickle request to the API.
func (c *IBKRClient) Tickle() error {
	resp, err := c.client.R().
		Post(fmt.Sprintf("%s/tickle", c.baseUrl))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to send tickle request")
	}

	return nil
}

// Reauthenticate reauthenticates the session with the Brokerage system.
func (c *IBKRClient) Reauthenticate() (bool, error) {
	type Response struct {
		Authenticated bool     `json:"authenticated"`
		Competing     bool     `json:"competing"`
		Connected     bool     `json:"connected"`
		Fail          string   `json:"fail"`
		Message       string   `json:"message"`
		Prompts       []string `json:"prompts"`
	}

	resp, err := c.client.R().
		Post(fmt.Sprintf("%s/iserver/reauthenticate", c.baseUrl))

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("failed to reauthenticate session")
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return false, err
	}

	return response.Connected, nil
}
