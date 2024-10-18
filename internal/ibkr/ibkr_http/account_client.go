package ibkr_http

import (
	"fmt"

	"encoding/json"
)

func (c *IBKRHttpClient) GetIServerAccounts() ([]string, error) {
	type Response struct {
		Accounts        []string          `json:"accounts"`
		Aliases         map[string]string `json:"aliases"`
		SelectedAccount string            `json:"selectedAccount"`
	}

	url := fmt.Sprintf("%s/iserver/accounts", c.url)

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get iServer accounts: %s", resp.Error())
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	fmt.Println(response)

	return response.Accounts, nil
}
