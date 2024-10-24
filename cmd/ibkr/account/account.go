package cmd

import (
	"os"

	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var GetIServerAccountsCmd = &cobra.Command{
	Use:   "getIserverAccounts",
	Short: "Get brokerage accounts",
	Long:  "Get brokerage accounts",
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		ibrkClient := ibkr.NewIBKRClient(logger)
		result, err := ibrkClient.HttpClient.GetIServerAccounts()
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		logger.Info().Msgf("i %s", result)

	},
}
