package cmd

import (
	"os"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var CheckAuthStatusCmd = &cobra.Command{
	Use:   "checkAuthStatus",
	Short: "Check SSO authentication status",
	Long:  "Check SSO authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		ibkrConfig := config.GetAppConfig().IbkrConfig
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		ibrkClient := ibkr.NewIBKRClient(ibkrConfig, logger)
		status, err := ibrkClient.HttpClient.GetAuthenticationStatus()
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		logger.Info().Msgf("Authentication status: %v", status)

	},
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from SSO",
	Long:  "Logout from SSO",
	Run: func(cmd *cobra.Command, args []string) {
		ibkrConfig := config.GetAppConfig().IbkrConfig
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		ibrkClient := ibkr.NewIBKRClient(ibkrConfig, logger)
		result, err := ibrkClient.HttpClient.Logout()
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		logger.Info().Msgf("Logout %s", result)

	},
}

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate SSO authentication",
	Long:  "Validate SSO authentication",
	Run: func(cmd *cobra.Command, args []string) {
		ibkrConfig := config.GetAppConfig().IbkrConfig
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		ibrkClient := ibkr.NewIBKRClient(ibkrConfig, logger)
		valid, err := ibrkClient.HttpClient.ValidateSession()
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		logger.Info().Msgf("Authentication is valid: %v", valid)

	},
}

var TickleCmd = &cobra.Command{
	Use:   "tickle",
	Short: "Tickle to ping the server to keep the session open",
	Long:  "Tickle to ping the server to keep the session open",
	Run: func(cmd *cobra.Command, args []string) {
		ibkrConfig := config.GetAppConfig().IbkrConfig
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		ibrkClient := ibkr.NewIBKRClient(ibkrConfig, logger)
		session, err := ibrkClient.HttpClient.Tickle()
		if err != nil {
			logger.Error().Msgf("%v", err)
			return
		}

		logger.Info().Msgf("session: %v", session)
	},
}
