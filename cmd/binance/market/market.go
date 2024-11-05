package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bincentive-ben/exchange"
	exbinance "github.com/bincentive-ben/exchange/binance"
	"github.com/btc-etf-arbitrage/internal/binance"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var IsWebsocket bool
var CurrencyList []string

var GetBalancesCmd = &cobra.Command{
	Use:   "getBlances",
	Short: "Get balances",
	Long:  "Get balances",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		cfg := config.GetAppConfig()
		client := binance.NewBinanceClient(logger, cfg.BinanceConfig.Type)
		balances, err := client.GetBalances()
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		for _, balance := range balances {
			for _, currency := range CurrencyList {
				if balance.Currency == currency {
					logger.Debug().Msgf("balance: %v", balance)
				}
			}
		}
	},
}

var CreateOrder = &cobra.Command{
	Use:   "createOrder",
	Short: "Create order",
	Long:  "Create order",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

		var quanity float64
		if len(args) == 0 {
			quanity = 0.00008
		} else {
			var err error
			quanity, err = strconv.ParseFloat(args[0], 64)
			fmt.Println("??")
			if err != nil {
				logger.Error().Msgf("%v", err)
				return
			}
		}

		cfg := config.GetAppConfig()
		logger.Debug().Msgf("quanity %v", quanity)
		client := binance.NewBinanceClient(logger, cfg.BinanceConfig.Type)

		param := map[string]interface{}{
			"SideEffect": exbinance.MarginBuy,
		}
		order, err := client.CreateOrder("BTCUSDT", exchange.OrderMarket, exchange.OrderSell, quanity, 0, param)
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		logger.Info().Msgf("order %v", order)
	},
}

var GetTickerPrice = &cobra.Command{
	Use:   "getTickerPrice",
	Short: "Get ticker price",
	Long:  "Get ticker price",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		cfg := config.GetAppConfig()
		client := binance.NewBinanceClient(logger, cfg.BinanceConfig.Type)
		if IsWebsocket {
			ctx, cancel := signal.NotifyContext(
				context.Background(),
				os.Interrupt,    // SIGINT
				syscall.SIGTERM, // SIGTERM
			)

			receiver := make(chan interface{}, 128)
			err := client.Exchange.Subscribe(exchange.Subscribe{
				Topic: exchange.TopicOrderBook,
				Param: map[string]interface{}{
					"Depth": 5,
				},
				Symbol: "BTCUSDT",
			}, receiver)
			if err != nil {
				return
			}

			go func() {
				for {
					select {
					case c := <-receiver:
						switch t := c.(type) {
						case exchange.OrderBook:
							logger.Debug().Msgf("OrderBook: bids %v", t.Bids)
							logger.Debug().Msgf("OrderBook: asks %v", t.Asks)
						}
					}
				}
			}()

			<-ctx.Done()
			cancel()
		} else {
			orderBook, err := client.Exchange.GetOrderBook(context.Background(), "BTCUSDT")
			if err != nil {
				logger.Error().Msgf("%v", err)
				return
			}

			logger.Info().Msgf("orderBook %v", orderBook)
		}

	},
}

func init() {
	GetTickerPrice.Flags().BoolVarP(&IsWebsocket, "websocket", "w", false, "Use websocket to get orderbook")
	GetBalancesCmd.Flags().StringArrayVarP(&CurrencyList, "list", "l", []string{"BTC", "USDT"}, "Currency list")
}
