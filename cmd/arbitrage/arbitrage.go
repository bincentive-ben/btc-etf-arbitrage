package cmd

import (
	binanceCmd "github.com/btc-etf-arbitrage/cmd/binance"
	ibkrCmd "github.com/btc-etf-arbitrage/cmd/ibkr"
	"github.com/btc-etf-arbitrage/internal/arbitrage"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "arbitrage",
	Short: "Start arbitrage application",
	Long:  "Start arbitrage application",
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run arbitrage application",
	Long:  "Run arbitrage application",
	Run: func(cmd *cobra.Command, args []string) {
		// Start the application
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}

		logger := arb.GetLogger()

		err = arb.StartAutoRefreshTraderConfig()
		if err != nil {
			logger.Error().Msgf("Error start auto refreshing trader config: %v", err)
			panic(err)
		}

		arb.Run()

	},
}

func init() {
	Cmd.AddCommand(
		ibkrCmd.IbkrCmd,
		binanceCmd.BinanceCmd,
		RunCmd,
	)
}
