package ibkr_http

import (
	"fmt"

	"encoding/json"
)

type Securities struct {
	Conid         string `json:"conid"`
	CompanyHeader string `json:"companyHeader"`
	CompanyName   string `json:"companyName"`
	Symbol        string `json:"symbol"`
	Description   string `json:"description"`
	Restricted    string `json:"restricted"`
	Fop           string `json:"fop"`
	Opt           string `json:"opt"`
	War           string `json:"war"`
	Sections      []struct {
		SecType    string `json:"secType"`
		Months     string `json:"months"`
		Symbol     string `json:"symbol"`
		Exchange   string `json:"exchange"`
		LegSecType string `json:"legSecType"`
	} `json:"sections"`
}

func (c *IBKRHttpClient) SearchIServerSecuriies(tickerList []string) error {

	for _, ticker := range tickerList {
		securities, err := c.SearchIServerSecurity(ticker)
		if err != nil {
			return err
		}
		for index, security := range securities {
			fmt.Printf("ticker: %s, index: %d, conid: %s companyHeader: %s, companyName: %s\n", ticker, index, security.Conid, security.CompanyHeader, security.CompanyName)
		}
	}

	return nil
}

// SearchIServerSecurity searches for security based on the provided query.
func (c *IBKRHttpClient) SearchIServerSecurity(query string) ([]Securities, error) {
	type Request struct {
		Symbol  string `json:"symbol"`
		Name    bool   `json:"name"`
		SecType string `json:"secType"`
	}

	url := fmt.Sprintf("%s/iserver/secdef/search", c.url)

	requestBody := Request{
		Symbol:  query,
		Name:    false,
		SecType: "",
	}

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to search securities: %s", resp.Error())
	}

	var securities []Securities
	err = json.Unmarshal(resp.Body(), &securities)
	if err != nil {
		return nil, err
	}

	return securities, nil
}
