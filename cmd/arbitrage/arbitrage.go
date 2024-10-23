package cmd

import (
	binanceCmd "github.com/btc-etf-arbitrage/cmd/binance"
	ibkrCmd "github.com/btc-etf-arbitrage/cmd/ibkr"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "arbitrage",
	Short: "Start arbitrage application",
	Long:  "Start arbitrage application",
}

func init() {
	Cmd.AddCommand(
		ibkrCmd.IbkrCmd,
		binanceCmd.BinanceCmd,
	)
}
