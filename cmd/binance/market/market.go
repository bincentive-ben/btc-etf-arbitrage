package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bincentive-ben/exchange"
	"github.com/btc-etf-arbitrage/internal/binance"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var IsWebsocket bool

var GetTickerPrice = &cobra.Command{
	Use:   "getTickerPrice",
	Short: "Get ticker price",
	Long:  "Get ticker price",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		client := binance.NewBinanceClient(logger)
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
}
