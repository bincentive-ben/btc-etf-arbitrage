package binance

import (
	"github.com/bincentive-ben/exchange"
	"github.com/bincentive-ben/exchange/binance"
	"github.com/bincentive-ben/exchange/common"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/rs/zerolog"
)

type BinanceClient struct {
	Exchange *common.AsyncExchange
	Config   *config.BinanceConfig
	logger   zerolog.Logger
}

// NewBinanceClient creates a new instance of BinanceClient
func NewBinanceClient(logger zerolog.Logger) *BinanceClient {
	cfg := config.GetAppConfig()
	binanceConfig := cfg.BinanceConfig

	account := exchange.ExchangeAccount{
		ExchangeID:   binanceConfig.ExchangeID,
		RESTEndpoint: binanceConfig.RestEndpoint,
		WSEndpoint:   binanceConfig.WsEndpoint,
		Apikey:       binanceConfig.ApiKey,
		Secret:       binanceConfig.SecretKey,
	}

	return &BinanceClient{
		Exchange: binance.NewSpotExchange(account),
		logger:   logger.With().Str("component", "binance").Logger(),
	}
}

func (c *BinanceClient) SubscribeExchange(receiver chan interface{}) error {
	c.logger.Debug().Msg("Start SubscribeExchange")

	subscribe := exchange.Subscribe{
		Topic:  exchange.TopicOrderBook,
		Symbol: "BTCUSDT",
		Param: map[string]interface{}{
			"Depth": 5,
		},
	}

	err := c.Exchange.Subscribe(subscribe, receiver)
	if err != nil {
		return err
	}

	return nil
}

func (b *BinanceClient) Subscribe(sub exchange.Subscribe, receiver chan interface{}) error {
	return b.Exchange.Subscribe(sub, receiver)
}
