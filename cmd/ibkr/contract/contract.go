package cmd

import (
	"github.com/btc-etf-arbitrage/internal/arbitrage"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/spf13/cobra"
)

var SearchIServerSecuritiesCmd = &cobra.Command{
	Use:   "searchIServerSecurities",
	Short: "Search securities of from 'StrategyConfig->StrategySettings->TickerList'",
	Long:  "Search securities of from 'StrategyConfig->StrategySettings->TickerList'",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		ibrkClient := ibkr.NewIBKRClient()
		tickerList := config.GetConfig().StrategyConfig.StrategySettings.TickerList
		err = ibrkClient.HttpClient.SearchIServerSecuriies(tickerList)
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

	},
}
