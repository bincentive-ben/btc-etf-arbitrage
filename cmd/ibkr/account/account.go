package cmd

import (
	"github.com/btc-etf-arbitrage/internal/arbitrage"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/spf13/cobra"
)

var GetIServerAccountsCmd = &cobra.Command{
	Use:   "getIserverAccounts",
	Short: "Get brokerage accounts",
	Long:  "Get brokerage accounts",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		// Brokerage accounts
		ibrkClient := ibkr.NewIBKRClient()
		result, err := ibrkClient.HttpClient.GetIServerAccounts()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		log.Info().Msgf("i %s", result)

	},
}
