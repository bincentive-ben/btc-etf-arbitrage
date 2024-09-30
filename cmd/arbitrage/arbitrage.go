package cmd

import (
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
	)
}
