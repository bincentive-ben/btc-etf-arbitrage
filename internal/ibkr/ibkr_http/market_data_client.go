package ibkr_http

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/btc-etf-arbitrage/internal/ibkr/common"
)

// GetIServerMarketDataSnapshot retrieves market data snapshot for the provided conids.
func (c *IBKRHttpClient) GetIServerMarketDataSnapshot(conids []string, since int, fields []string) ([]common.MarketDataSnapshot, error) {
	url := fmt.Sprintf("%s/iserver/marketdata/snapshot", c.url)

	conidsStr := strings.Join(conids, ",")
	fieldsStr := strings.Join(fields, ",")

	resp, err := c.client.R().
		SetQueryParam("conids", conidsStr).
		SetQueryParam("since", strconv.Itoa(since)).
		SetQueryParam("fields", fieldsStr).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to retrieve market data snapshot: %s", resp.Error())
	}

	var response []common.MarketDataSnapshot
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Not done yet
func (c *IBKRHttpClient) GetIServerMarketDataHistory(conid string, period string, barSize string) ([]common.MarketDataHistory, error) {
	url := fmt.Sprintf("%s/iserver/marketdata/history", c.url)

	resp, err := c.client.R().
		SetQueryParam("conid", conid).
		SetQueryParam("period", period).
		SetQueryParam("barSize", barSize).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to retrieve market data history: %s", resp.Error())
	}

	var response []common.MarketDataHistory
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
