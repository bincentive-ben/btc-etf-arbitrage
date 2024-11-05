package cmd

import (
	marketDataCmd "github.com/btc-etf-arbitrage/cmd/binance/market"
	"github.com/spf13/cobra"
)

var BinanceCmd = &cobra.Command{
	Use:   "binance",
	Short: "Binance commands",
	Long:  "Binance commands",
}

func init() {
	BinanceCmd.AddCommand(marketDataCmd.GetTickerPrice)
	BinanceCmd.AddCommand(marketDataCmd.CreateOrder)
	BinanceCmd.AddCommand(marketDataCmd.GetBalancesCmd)
}
