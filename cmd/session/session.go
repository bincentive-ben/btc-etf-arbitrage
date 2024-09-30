package cmd

import (
	"github.com/btc-etf-arbitrage/internal/arbitrage"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/spf13/cobra"
)

var CheckAuthStatusCmd = &cobra.Command{
	Use:   "checkAuthStatus",
	Short: "Check SSO authentication status",
	Long:  "Check SSO authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		// Get authentication status
		ibrkClient := ibkr.NewIBKRClient()
		status, err := ibrkClient.GetAuthenticationStatus()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		log.Info().Msgf("Authentication status: %v", status)

	},
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from SSO",
	Long:  "Logout from SSO",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		// Logout
		ibrkClient := ibkr.NewIBKRClient()
		result, err := ibrkClient.Logout()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		log.Info().Msgf("Logout %s", result)

	},
}
