package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btc-etf-arbitrage/internal/arbitrage"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/btc-etf-arbitrage/internal/ibkr/common"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_websocket"
	"github.com/spf13/cobra"
)

var IsWebsocket bool

var GetIServerMarketDataSnapshotCmd = &cobra.Command{
	Use:   "getIServerMarketDataSnapshotCmd",
	Short: "Get market data snapshot",
	Long:  "Get market data snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		contractIDs := config.GetConfig().StrategyConfig.StrategySettings.ContractIDList
		fields := config.GetConfig().AppConfig.IBKR.Fields
		var marketDataSnapshot []common.MarketDataSnapshot
		ibkrClient := ibkr.NewIBKRClient()

		if IsWebsocket {
			defer ibkrClient.WsClient.Close()
			go ibkrClient.WsClient.StartListening()

			select {
			case <-ibkrClient.WsClient.Authenticated:
				fmt.Println("Received authenticated signal, sending request for streaming data.")
				request := ibkr_websocket.StreamingDataRequest{
					Conid:  "677037673",
					Fields: []string{"31", "55", "84", "86"},
				}
				ibkrClient.WsClient.SubscribeStreamingData(request)
			case <-time.After(30 * time.Second):
				fmt.Println("Timeout waiting for authentication")
			}
		} else {
			marketDataSnapshot, err = ibkrClient.HttpClient.GetIServerMarketDataSnapshot(contractIDs, 0, fields)
			if err != nil {
				log.Error().Msgf("%v", err)
				return
			}

			for _, marketData := range marketDataSnapshot {
				fmt.Printf("Symbol: %s, LastPrice: %s BidPrice: %s, AskPrice: %s\n", marketData.Symbol, marketData.LastPrice, marketData.BidPrice, marketData.AskPrice)
			}
		}

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
	},
}

var GetIServerMarketDataHistoryCmd = &cobra.Command{
	Use:   "getIServerMarketDataHistoryCmd",
	Short: "Get market data history",
	Long:  "Get market data history",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		ibkrClient := ibkr.NewIBKRClient()
		fmt.Println("ibkrClient:", ibkrClient)
		if IsWebsocket {
			defer ibkrClient.WsClient.Close()
			go ibkrClient.WsClient.StartListening()

			select {
			case <-ibkrClient.WsClient.Authenticated:
				fmt.Println("Received authenticated signal, sending request for historical data.")
				request := ibkr_websocket.HistoricalDataRequest{
					Conid:  "677037676",
					Period: "1d",
					Bar:    "1hour",
					Source: "trades",
					Format: "%o/%c/%h/%l",
				}
				ibkrClient.WsClient.SubscribeHistoricalData(request)
			case <-time.After(30 * time.Second):
				fmt.Println("Timeout waiting for authentication")
			}

		} else {
			log.Debug().Msgf("Not done yet")
			// marketDataHistory, err = ibkrClient.HttpClient.GetIServerMarketDataHistory(contractIDs, 0, fields)
			// if err != nil {
			// 	log.Error().Msgf("%v", err)
			// 	return
			// }
		}

		// for _, marketData := range marketDataHistory {
		// 	fmt.Printf("Symbol: %s, LastPrice: %s BidPrice: %s, AskPrice: %s\n", marketData.Symbol, marketData.LastPrice, marketData.BidPrice, marketData.AskPrice)
		// }

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

	},
}

func init() {
	GetIServerMarketDataSnapshotCmd.Flags().BoolVarP(&IsWebsocket, "websocket", "w", false, "Use websocket to get market data snapshot")
	GetIServerMarketDataHistoryCmd.Flags().BoolVarP(&IsWebsocket, "websocket", "w", false, "Use websocket to get market data history")
}
