package binance

import (
	"context"
	"strings"

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

var BinanceReceiver = make(chan interface{}, 128)

func GetBinanceReceiver() chan interface{} {
	return BinanceReceiver
}

// NewBinanceClient creates a new instance of BinanceClient
func NewBinanceClient(logger zerolog.Logger, exchangeType string) *BinanceClient {
	cfg := config.GetAppConfig()
	binanceConfig := cfg.BinanceConfig

	account := exchange.ExchangeAccount{
		ExchangeID:   binanceConfig.ExchangeID,
		RESTEndpoint: binanceConfig.RestEndpoint,
		WSEndpoint:   binanceConfig.WsEndpoint,
		Apikey:       binanceConfig.ApiKey,
		Secret:       binanceConfig.SecretKey,
	}

	exchangeType = strings.ToLower(exchangeType)

	var exchange *common.AsyncExchange
	switch exchangeType {
	case "spot":
		exchange = binance.NewSpotExchange(account)
	case "margin":
		exchange = binance.NewMarginExchange(account)
	default:
		// use spot as default
		exchange = binance.NewSpotExchange(account)
	}

	return &BinanceClient{
		Exchange: exchange,
		logger:   logger.With().Str("component", "binance").Logger(),
	}
}

func (c *BinanceClient) SubscribeExchange(receiver chan interface{}) error {
	c.logger.Debug().Msg("Start SubscribeExchange")

	subscribe := exchange.Subscribe{
		Topic:  exchange.TopicOrderBook,
		Symbol: "BTCUSDT",
		Param: map[string]interface{}{
			"Interval": "3",
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

func (b *BinanceClient) ProcessBinanceMessage(receiver chan interface{}) {
	for c := range receiver {
		switch t := c.(type) {
		case exchange.OrderBook:
			b.logger.Debug().Msgf("binance bids: len:%v  first: %v", len(t.Bids), t.Bids[0].Price)
			b.logger.Debug().Msgf("binance Asks: len:%v  first: %v", len(t.Asks), t.Asks[0].Price)
			// a.logger.Debug().Msgf("binance asks: len:%v  %v", len(t.Asks), t.Asks)
		}
	}
}

func (b *BinanceClient) GetBalances() ([]exchange.Balance, error) {

	balances, err := b.Exchange.GetBalances(context.Background())
	if err != nil {
		b.logger.Error().Msgf("GetBalances error: %v", err)
	}

	// for _, balance := range balances {
	// 	if balance.Currency == "USDT" {
	// 		b.logger.Debug().Msgf("balance: %v", balance)
	// 	}
	// }

	return balances, nil
}

func (b *BinanceClient) CreateOrder(symbol string, orderType exchange.OrderType, orderSide exchange.OrderSide, quanity float64, price float64, param map[string]interface{}) (exchange.Order, error) {
	order, err := b.Exchange.CreateOrder(context.Background(), symbol, orderType, orderSide, quanity, price, param)
	if err != nil {
		b.logger.Error().Msgf("CreateOrder error: %v", err)
		return exchange.Order{}, err
	}

	b.logger.Debug().Msgf("CreateOrder: %v", order)
	return order, nil
}
