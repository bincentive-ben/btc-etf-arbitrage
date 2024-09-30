package cmd

import (
	"fmt"
	"os"

	arbitrageCmd "github.com/btc-etf-arbitrage/cmd/arbitrage"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/spf13/cobra"
)

var configFile string

var RootCmd = &cobra.Command{
	Use:   "btc-etf-arbitrage",
	Short: "Application to find arbitrage opportunities between BTC and ETFs",
	Long:  "Application to find arbitrage opportunities between BTC and ETFs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration
		err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "configs/default.yaml", "config file (default is configs/default.yaml)")
	RootCmd.AddCommand(
		arbitrageCmd.Cmd,
	)
}
