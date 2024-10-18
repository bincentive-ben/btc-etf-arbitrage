package arbitrage

import (
	"os"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/rs/zerolog"
)

type Arbitrage struct {
	cfg    config.Config
	Logger zerolog.Logger
}

func NewArbitrage() (*Arbitrage, error) {
	cfg := config.GetConfig()

	return &Arbitrage{
		cfg:    cfg,
		Logger: zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger(),
	}, nil
}
