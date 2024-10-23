package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bincentive-ben/exchange"
	"github.com/btc-etf-arbitrage/internal/arbitrage"
	"github.com/btc-etf-arbitrage/internal/binance"
	"github.com/spf13/cobra"
)

var IsWebsocket bool

var GetTickerPrice = &cobra.Command{
	Use:   "getTickerPrice",
	Short: "Get ticker price",
	Long:  "Get ticker price",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		client := binance.NewBinanceClient()
		if IsWebsocket {
			ctx, cancel := signal.NotifyContext(
				context.Background(),
				os.Interrupt,    // SIGINT
				syscall.SIGTERM, // SIGTERM
			)

			receiver := make(chan interface{}, 128)
			err = client.SpotExchange.Subscribe(exchange.Subscribe{
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
							log.Debug().Msgf("OrderBook: %v", t)
						}
					}
				}
			}()

			<-ctx.Done()
			cancel()
		} else {
			orderBook, err := client.SpotExchange.GetOrderBook(context.Background(), "BTCUSDT")
			if err != nil {
				log.Error().Msgf("%v", err)
				return
			}

			log.Info().Msgf("orderBook %v", orderBook)
		}

	},
}

func init() {
	GetTickerPrice.Flags().BoolVarP(&IsWebsocket, "websocket", "w", false, "Use websocket to get orderbook")
}
