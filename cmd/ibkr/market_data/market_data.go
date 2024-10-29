package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/btc-etf-arbitrage/internal/ibkr/common"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_websocket"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var IsWebsocket bool

var GetIServerMarketDataSnapshotCmd = &cobra.Command{
	Use:   "getIServerMarketDataSnapshotCmd",
	Short: "Get market data snapshot",
	Long:  "Get market data snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		ibkrConfig := config.GetAppConfig().IbkrConfig
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

		contractIDs := config.GetAppConfig().IbkrConfig.ContractIDList
		fields := config.GetAppConfig().IbkrConfig.Fields
		var marketDataSnapshot []common.MarketDataSnapshot
		var err error
		ibkrClient := ibkr.NewIBKRClient(ibkrConfig, logger)

		if IsWebsocket {
			ibkrReceiver := make(chan interface{}, 128)
			defer ibkrClient.WsClient.Close()
			go ibkrClient.WsClient.StartListening(ibkrReceiver)

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
				logger.Error().Msgf("%v", err)
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
		ibkrConfig := config.GetAppConfig().IbkrConfig
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

		ibkrClient := ibkr.NewIBKRClient(ibkrConfig, logger)
		fmt.Println("ibkrClient:", ibkrClient)
		if IsWebsocket {
			defer ibkrClient.WsClient.Close()
			ibkrReceiver := make(chan interface{}, 128)
			go ibkrClient.WsClient.StartListening(ibkrReceiver)
			go ibkrClient.ProcessMessage(ibkrReceiver)

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
			logger.Debug().Msgf("Not done yet")
			// marketDataHistory, err = ibkrClient.HttpClient.GetIServerMarketDataHistory(contractIDs, 0, fields)
			// if err != nil {
			// 	logger.Error().Msgf("%v", err)
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
