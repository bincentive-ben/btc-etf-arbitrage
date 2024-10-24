package cmd

import (
	"fmt"
	"os"

	arbitrageCmd "github.com/btc-etf-arbitrage/cmd/arbitrage"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "btc-etf-arbitrage",
	Short: "Application to find arbitrage opportunities between BTC and ETFs",
	Long:  "Application to find arbitrage opportunities between BTC and ETFs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration
		err := config.LoadAppConfig(config.ConfigFile)
		if err != nil {
			fmt.Println("Error loading app config:", err)
			os.Exit(1)
		}
		err = config.LoadTraderConfig(config.TraderConfigFile)
		if err != nil {
			fmt.Println("Error loading trader config:", err)
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
	RootCmd.PersistentFlags().StringVar(&config.ConfigFile, "config", "configs/default.yaml", "app config file (default is configs/default.yaml)")
	RootCmd.PersistentFlags().StringVar(&config.TraderConfigFile, "traderConfig", "configs/trader_default.yaml", "trader config file (default is configs/trader_default.yaml)")
	RootCmd.AddCommand(
		arbitrageCmd.Cmd,
	)
}
