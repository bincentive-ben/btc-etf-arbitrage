package arbitrage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bincentive-ben/exchange"
	exbinance "github.com/bincentive-ben/exchange/binance"
	"github.com/btc-etf-arbitrage/internal/binance"
	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_websocket"

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

	logger := zerolog.New(multi).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	// logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ibkrClient := ibkr.NewIBKRClient(appConfig.IbkrConfig, logger)
	ibkrClient.StartScheduler()

	return &Arbitrage{
		appConfig:     appConfig,
		traderConfig:  traderConfig,
		ibkrClient:    ibkrClient,
		binanceClient: binance.NewBinanceClient(logger, appConfig.BinanceConfig.Type),
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

	err := a.ProcessMessage()
	if err != nil {
		a.logger.Error().Msgf("Error process message: %v", err)
		return err
	}

	err = a.SubscribeBinanceExchange()
	if err != nil {
		a.logger.Error().Msgf("Error subscribe binance exchange: %v", err)
		return err
	}

	err = a.SubscribeIbkrExchange()
	if err != nil {
		a.logger.Error().Msgf("Error subscribe ibkr exchange: %v", err)
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

func (a *Arbitrage) SubscribeBinanceExchange() error {
	a.logger.Debug().Msg("Start SubscribeBinanceExchange")

	receiver := binance.GetBinanceReceiver()
	client := a.GetBinanceClient()
	err := client.SubscribeExchange(receiver)
	if err != nil {
		return err
	}

	return nil
}

func (a *Arbitrage) SubscribeIbkrExchange() error {
	a.logger.Debug().Msg("Start SubscribeIbkrExchange")

	receiver := ibkr.GetIbkrReceiver()
	client := a.GetIbkrClient()
	go client.StartListening(receiver)

	err := client.SubscribeExchange(receiver)
	if err != nil {
		return err
	}

	return nil
}

func (a *Arbitrage) ProcessMessage() error {
	ibkrReceiver := ibkr.GetIbkrReceiver()
	binanceReceiver := binance.GetBinanceReceiver()

	go a.ProcessIbkrMessage(ibkrReceiver)
	go a.ProcessBinanceMessage(binanceReceiver)

	return nil
}

func (a *Arbitrage) ProcessIbkrMessage(receiver chan interface{}) error {
	for m := range receiver {
		message, ok := m.([]byte)
		if !ok {
			a.logger.Error().Msg("Error type asserting message")
			return fmt.Errorf("error type asserting message")
		}

		a.logger.Debug().Msgf("received message: %s", string(message))

		var msg ibkr_websocket.Message
		err := json.Unmarshal(message, &msg)
		if err != nil {
			a.logger.Error().Msgf("Error unmarshalling message: %v", err)
			continue
		}

		switch msg.Topic {
		case "sts":
			var stsMsg ibkr_websocket.StsMessage
			if err := json.Unmarshal(message, &stsMsg); err != nil {
				a.logger.Error().Msgf("Error unmarshalling message: %v", err)
				continue
			}

			// Check authenticated flag
			if stsMsg.Topic == "sts" && stsMsg.Args.Authenticated {
				a.logger.Debug().Msg("Authenticated!")

				ibkrClient := a.GetIbkrClient()
				ibkrClient.AuthenticatedCh <- true
			}

		case "sor":
			var sorMsg ibkr_websocket.SorMessage
			if err := json.Unmarshal(message, &sorMsg); err != nil {
				a.logger.Error().Msgf("Error unmarshalling message: %v", err)
				continue
			}

			for index, arg := range sorMsg.Args {
				a.logger.Debug().Msgf("sorMsg.Args[%d]: %v", index, arg)
				if hasOrderFilled(arg) {
					a.logger.Debug().Msgf("Order filled! %v, %v, %v, %v", arg.OrderID, arg.Conid, arg.Price, arg.FilledQuantity)
					a.logger.Debug().Msgf("Start to hedge on binance")
					binanceClient := a.GetBinanceClient()

					param := map[string]interface{}{
						"SideEffect": exbinance.MarginBuy,
					}

					order, err := binanceClient.CreateOrder("BTCUSDT", exchange.OrderMarket, exchange.OrderSell, 0.00008, 0, param)
					if err != nil {
						a.logger.Error().Msgf("%v", err)
						return err
					}

					a.logger.Info().Msgf("order %v", order)
				}
			}

		case "sbd":
		case "str":
		case "spl":
		case "act":
		default:
		}

	}
	return nil
}

func (a *Arbitrage) ProcessBinanceMessage(receiver chan interface{}) {
	for c := range receiver {
		switch t := c.(type) {
		case exchange.OrderBook:
			a.logger.Debug().Msgf("binance bids: len:%v  first: %v", len(t.Bids), t.Bids[0].Price)
			a.logger.Debug().Msgf("binance Asks: len:%v  first: %v", len(t.Asks), t.Asks[0].Price)
			// a.logger.Debug().Msgf("binance asks: len:%v  %v", len(t.Asks), t.Asks)
		default:
			// TODO: handle other types
		}
	}
}

func RefreshTraderConfig() error {
	err := config.LoadTraderConfig(config.TraderConfigFile)
	if err != nil {
		return err
	}

	return nil
}

func hasOrderFilled(sorArgs ibkr_websocket.SorArgs) bool {
	return sorArgs.Status == "Filled"
}
