package arbitrage

import (
	"os"

	"github.com/btc-etf-arbitrage/internal/binance"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/rs/zerolog"
)

type Arbitrage struct {
	cfg           config.Config
	IbkrClient    *ibkr.IBKRClient
	BinanceClient *binance.BinanceClient
	Logger        zerolog.Logger
}

func NewArbitrage() (*Arbitrage, error) {
	cfg := config.GetConfig()

	return &Arbitrage{
		cfg:           cfg,
		IbkrClient:    ibkr.NewIBKRClient(),
		BinanceClient: binance.NewBinanceClient(),
		Logger:        zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger(),
	}, nil
}
