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
		status, err := ibrkClient.HttpClient.GetAuthenticationStatus()
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
		result, err := ibrkClient.HttpClient.Logout()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		log.Info().Msgf("Logout %s", result)

	},
}

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate SSO authentication",
	Long:  "Validate SSO authentication",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		// Validate authentication
		ibrkClient := ibkr.NewIBKRClient()
		valid, err := ibrkClient.HttpClient.ValidateSession()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		log.Info().Msgf("Authentication is valid: %v", valid)

	},
}

var TickleCmd = &cobra.Command{
	Use:   "tickle",
	Short: "Tickle to ping the server to keep the session open",
	Long:  "Tickle to ping the server to keep the session open",
	Run: func(cmd *cobra.Command, args []string) {
		arb, err := arbitrage.NewArbitrage()
		if err != nil {
			panic(err)
		}
		log := arb.Logger

		// Tickle
		ibrkClient := ibkr.NewIBKRClient()
		session, err := ibrkClient.HttpClient.Tickle()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		log.Info().Msgf("session: %v", session)
	},
}
