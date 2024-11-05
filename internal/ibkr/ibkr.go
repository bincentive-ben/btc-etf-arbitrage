package ibkr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_http"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_websocket"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/go-co-op/gocron"
)

type IbkrClient struct {
	HttpClient      *ibkr_http.IBKRHttpClient
	WsClient        *ibkr_websocket.IBKRWebsocketClient
	Scheduler       *gocron.Scheduler
	Config          config.IbkrConfig
	logger          zerolog.Logger
	AuthenticatedCh chan bool
}

var IbkrReceiver = make(chan interface{}, 128)

func GetIbkrReceiver() chan interface{} {
	return IbkrReceiver
}

// NewIBKRClient creates a new instance of IbkrClient
func NewIBKRClient(config config.IbkrConfig, logger zerolog.Logger) *IbkrClient {
	httpClient := ibkr_http.NewIBKRHttpClient()
	wsClient, err := ibkr_websocket.NewIBKRWebsocketClient(config.WsEndpoint, logger)
	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	return &IbkrClient{
		HttpClient:      httpClient,
		WsClient:        wsClient,
		Scheduler:       scheduler,
		Config:          config,
		AuthenticatedCh: make(chan bool),
		logger:          logger.With().Str("component", "ibkr").Logger(),
	}
}

func (c IbkrClient) GetScheduler() *gocron.Scheduler {
	return c.Scheduler
}

func (c IbkrClient) StartScheduler() {
	scheduler := c.Scheduler

	// Ping session to keep the session alive every minute
	scheduler.Every(1).Minutes().Do(c.WsClient.PingSession)
	scheduler.Every(1).Minutes().Do(c.HttpClient.Tickle())

	scheduler.StartAsync()

}

func (c IbkrClient) SubscribeExchange(receiver chan interface{}) error {
	c.logger.Debug().Msg("Start SubscribeExchange")
	authenticated := false

	// Subscribe exchange after receiving authenticated signal
	select {
	case <-c.AuthenticatedCh:
		authenticated = true
	case <-time.After(5 * time.Second):
	}

	if !authenticated {
		return fmt.Errorf("timeout waiting for authentication")
	}

	// c.SubscribeStreamingDataList(c.Config.ContractIDList)
	// c.SubscribeHistoricalData("677037676")
	c.SubscribeLiveOrderUpdate()

	return nil
}

func (c IbkrClient) SubscribeStreamingDataList(conids []string) error {
	for _, conid := range conids {
		request := ibkr_websocket.StreamingDataRequest{
			Conid:  conid,
			Fields: c.Config.Fields,
		}
		err := c.WsClient.SubscribeStreamingData(request)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c IbkrClient) SubscribeStreamingData(request ibkr_websocket.StreamingDataRequest) error {
	err := c.WsClient.SubscribeStreamingData(request)
	if err != nil {
		return err
	}

	return nil
}

func (c IbkrClient) SubscribeHistoricalData(conid string) error {
	request := ibkr_websocket.HistoricalDataRequest{
		Conid:  conid,
		Period: "1d",
		Bar:    "1hour",
		Source: "trades",
		Format: "%o/%c/%h/%l",
	}

	err := c.WsClient.SubscribeHistoricalData(request)
	if err != nil {
		return err
	}

	return nil
}

func (c IbkrClient) SubscribeLiveOrderUpdate() error {
	err := c.WsClient.SubscribeLiveOrderUpdate()
	if err != nil {
		return err
	}

	return nil
}

func (c IbkrClient) GetAuthenticationStatus() (ibkr_http.GetAuthenticationStatusResponse, error) {
	authenticationStatus, err := c.HttpClient.GetAuthenticationStatus()
	if err != nil {
		return ibkr_http.GetAuthenticationStatusResponse{}, err
	}

	c.logger.Debug().Msgf("Authentication status: %v", authenticationStatus)

	return authenticationStatus, err
}

func (c IbkrClient) StartListening(receiver chan interface{}) error {
	conn, err := c.WsClient.GetConn()
	if err != nil {
		return err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			c.logger.Error().Msgf("Error reading message: %v", err)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.logger.Debug().Msg("Connection closed normally, attempting to reconnect...")

				c.logger.Debug().Msg("Start reauthenticating")
				authenticationStatus, err := c.GetAuthenticationStatus()
				if err != nil {
					c.logger.Error().Msgf("Error getting authentication status: %v", err)
				}
				c.logger.Debug().Msgf("Authentication status: %v", authenticationStatus.Authenticated)
				if !authenticationStatus.Authenticated {
					c.logger.Debug().Msg("Reauthenticating...")
					res, err := c.HttpClient.Reauthenticate()
					if err != nil {
						c.logger.Error().Msgf("Error reauthenticating: %v", err)
					}
					c.logger.Debug().Msgf("Reauthentication response: %v", res)
				}

				if err = c.WsClient.Reconnect(); err != nil {
					c.logger.Error().Msgf("Error reconnecting to websocket: %v", err)
					return err
				}
			} else {
				return err
			}
		}

		receiver <- message
	}
}

func (c IbkrClient) ProcessIbkrMessage(receiver chan interface{}) error {
	for m := range receiver {
		message, ok := m.([]byte)
		if !ok {
			c.logger.Error().Msg("Error type asserting message")
			return fmt.Errorf("error type asserting message")
		}

		c.logger.Debug().Msgf("received message: %s", string(message))

		var msg ibkr_websocket.Message
		err := json.Unmarshal(message, &msg)
		if err != nil {
			c.logger.Error().Msgf("Error unmarshalling message: %v", err)
			continue
		}

		switch msg.Topic {
		case "sts":
			var stsMsg ibkr_websocket.StsMessage
			if err := json.Unmarshal(message, &stsMsg); err != nil {
				c.logger.Error().Msgf("Error unmarshalling message: %v", err)
				continue
			}

			// Check authenticated flag
			if stsMsg.Topic == "sts" && stsMsg.Args.Authenticated {
				c.logger.Debug().Msg("Authenticated!")
				c.AuthenticatedCh <- true
			}

		case "sor":
			var sorMsg ibkr_websocket.SorMessage
			if err := json.Unmarshal(message, &sorMsg); err != nil {
				c.logger.Error().Msgf("Error unmarshalling message: %v", err)
				continue
			}

			for index, arg := range sorMsg.Args {
				c.logger.Debug().Msgf("sorMsg.Args[%d]: %v", index, arg)
				if hasOrderFilled(arg) {
					c.logger.Debug().Msgf("Order filled! %v, %v, %v, %v", arg.OrderID, arg.Conid, arg.Price, arg.FilledQuantity)
					c.logger.Debug().Msgf("Going to hedge on binance")
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

func hasOrderFilled(sorArgs ibkr_websocket.SorArgs) bool {
	return sorArgs.Status == "Filled"
}
