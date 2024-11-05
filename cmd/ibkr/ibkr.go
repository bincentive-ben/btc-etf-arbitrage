package cmd

import (
	accountCmd "github.com/btc-etf-arbitrage/cmd/ibkr/account"
	contractCmd "github.com/btc-etf-arbitrage/cmd/ibkr/contract"
	marketDataCmd "github.com/btc-etf-arbitrage/cmd/ibkr/market_data"
	sessionCmd "github.com/btc-etf-arbitrage/cmd/ibkr/session"

	"github.com/spf13/cobra"
)

var IbkrCmd = &cobra.Command{
	Use:   "ibkr",
	Short: "IBKR commands",
	Long:  "IBKR commands",
}

func init() {
	IbkrCmd.AddCommand(
		sessionCmd.CheckAuthStatusCmd,
		sessionCmd.LogoutCmd,
		sessionCmd.ValidateCmd,
		sessionCmd.TickleCmd,
		sessionCmd.ReauthenticateCmd,
		accountCmd.GetIServerAccountsCmd,
		contractCmd.SearchIServerSecuritiesCmd,
		marketDataCmd.GetIServerMarketDataSnapshotCmd,
		marketDataCmd.GetIServerMarketDataHistoryCmd,
	)
}
