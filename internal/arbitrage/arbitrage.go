package arbitrage

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btc-etf-arbitrage/internal/binance"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
)

type Arbitrage struct {
	appConfig     config.AppConfig
	traderConfig  config.TraderConfig
	ibkrClient    *ibkr.IbkrClient
	binanceClient *binance.BinanceClient
	scheduler     *gocron.Scheduler
	logger        zerolog.Logger
}

func NewArbitrage() (*Arbitrage, error) {
	appConfig := config.GetAppConfig()
	traderConfig := config.GetTraderConfig()
	scheduler := gocron.NewScheduler(time.UTC)

	logFile, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

	logger := zerolog.New(multi).Level(zerolog.DebugLevel).With().Time("time", time.Now()).Logger()
	// logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	return &Arbitrage{
		appConfig:     appConfig,
		traderConfig:  traderConfig,
		ibkrClient:    ibkr.NewIBKRClient(appConfig.IbkrConfig, logger),
		binanceClient: binance.NewBinanceClient(logger),
		scheduler:     scheduler,
		logger:        logger.With().Str("app", "arbitrage").Logger(),
	}, nil
}

func (a *Arbitrage) Run() error {
	a.logger.Debug().Msg("Start running arbitrage")
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,    // SIGINT
		syscall.SIGTERM, // SIGTERM
	)

	ibkrReceiver := make(chan interface{}, 128)
	binanceReceiver := make(chan interface{}, 128)

	err := a.ProcessMessage(ibkrReceiver, binanceReceiver)
	if err != nil {
		a.logger.Err(err).Msg("Error process message")
		return err
	}

	// err = a.SubscribeBinanceExchange(binanceReceiver)
	// if err != nil {
	// 	a.logger.Err(err).Msg("Error subscribe binance exchange")
	// 	return err
	// }

	err = a.SubscribeIbkrExchange(ibkrReceiver)
	if err != nil {
		a.logger.Err(err).Msg("Error subscribe ibkr exchange")
		return err
	}

	<-ctx.Done()
	a.logger.Debug().Msg("End running arbitrage")
	cancel()

	return nil
}

func (a *Arbitrage) GetAppConfig() config.AppConfig {
	return a.appConfig
}

func (a *Arbitrage) GetTraderConfig() config.TraderConfig {
	return a.traderConfig
}

func (a *Arbitrage) GetIbkrClient() *ibkr.IbkrClient {
	return a.ibkrClient
}

func (a *Arbitrage) GetBinanceClient() *binance.BinanceClient {
	return a.binanceClient
}

func (a *Arbitrage) GetScheduler() *gocron.Scheduler {
	return a.scheduler
}

func (a *Arbitrage) GetLogger() zerolog.Logger {
	return a.logger
}

func (a *Arbitrage) StartAutoRefreshTraderConfig() error {
	seconds := a.GetAppConfig().ReloadTraderConfigSeconds
	scheduler := a.GetScheduler()

	_, err := a.scheduler.Every(seconds).Seconds().Do(RefreshTraderConfig)
	if err != nil {
		return err
	}

	scheduler.StartAsync()

	return nil
}

func (a *Arbitrage) SubscribeBinanceExchange(receiver chan interface{}) error {
	a.logger.Debug().Msg("Start SubscribeBinanceExchange")
	client := a.GetBinanceClient()
	err := client.SubscribeExchange(receiver)
	if err != nil {
		return err
	}

	return nil
}

func (a *Arbitrage) SubscribeIbkrExchange(receiver chan interface{}) error {
	a.logger.Debug().Msg("Start SubscribeIbkrExchange")

	client := a.GetIbkrClient()
	go client.StartListening(receiver)

	err := client.SubscribeExchange(receiver)
	if err != nil {
		return err
	}

	return nil
}

func (a *Arbitrage) ProcessMessage(ibkrReceiver, binanceReceiver chan interface{}) error {
	go a.ibkrClient.ProcessMessage(ibkrReceiver)
	go a.binanceClient.ProcessMessage(binanceReceiver)

	return nil
}

func RefreshTraderConfig() error {
	err := config.LoadTraderConfig(config.TraderConfigFile)
	if err != nil {
		return err
	}

	return nil
}
