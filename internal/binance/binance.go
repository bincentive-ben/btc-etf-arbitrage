package binance

import (
	"github.com/bincentive-ben/exchange"
	"github.com/bincentive-ben/exchange/binance"
	"github.com/bincentive-ben/exchange/common"
	"github.com/btc-etf-arbitrage/internal/config"
)

type BinanceClient struct {
	SpotExchange *common.AsyncExchange
	Config       *config.Config
}

// NewBinanceClient creates a new instance of BinanceClient
func NewBinanceClient() *BinanceClient {
	cfg := config.GetConfig()
	binanceConfig := cfg.AppConfig.Binance

	account := exchange.ExchangeAccount{
		ExchangeID:   binanceConfig.ExchangeID,
		RESTEndpoint: binanceConfig.RestEndpoint,
		WSEndpoint:   binanceConfig.WsEndpoint,
		Apikey:       binanceConfig.ApiKey,
		Secret:       binanceConfig.SecretKey,
	}

	return &BinanceClient{
		SpotExchange: binance.NewSpotExchange(account),
	}
}
