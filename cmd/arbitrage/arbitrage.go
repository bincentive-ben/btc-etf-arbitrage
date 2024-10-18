package cmd

import (
	accountCmd "github.com/btc-etf-arbitrage/cmd/account"
	contractCmd "github.com/btc-etf-arbitrage/cmd/contract"
	marketDataCmd "github.com/btc-etf-arbitrage/cmd/market_data"
	sessionCmd "github.com/btc-etf-arbitrage/cmd/session"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "arbitrage",
	Short: "Start arbitrage application",
	Long:  "Start arbitrage application",
}

func init() {
	Cmd.AddCommand(
		sessionCmd.CheckAuthStatusCmd,
		sessionCmd.LogoutCmd,
		sessionCmd.ValidateCmd,
		sessionCmd.TickleCmd,
		accountCmd.GetIServerAccountsCmd,
		contractCmd.SearchIServerSecuritiesCmd,
		marketDataCmd.GetIServerMarketDataSnapshotCmd,
		marketDataCmd.GetIServerMarketDataHistoryCmd,
	)
}
