package cmd

import (
	"os"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var SearchIServerSecuritiesCmd = &cobra.Command{
	Use:   "searchIServerSecurities",
	Short: "Search securities of from 'StrategyConfig->StrategySettings->TickerList'",
	Long:  "Search securities of from 'StrategyConfig->StrategySettings->TickerList'",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

		ibrkClient := ibkr.NewIBKRClient(logger)
		tickerList := config.GetAppConfig().IbkrConfig.TickerList

		err := ibrkClient.HttpClient.SearchIServerSecuriies(tickerList)
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

	},
}
