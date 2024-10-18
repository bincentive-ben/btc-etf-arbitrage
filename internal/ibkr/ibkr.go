package ibkr

import (
	"time"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_http"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_websocket"

	"github.com/go-co-op/gocron"
)

type IBKRClient struct {
	HttpClient *ibkr_http.IBKRHttpClient
	WsClient   *ibkr_websocket.IBKRWebsocketClient
	Scheduler  *gocron.Scheduler
	Config     *config.Config

	AuthenticatedCh chan bool
}

// NewIBKRClient creates a new instance of IBKRClient
func NewIBKRClient() *IBKRClient {
	ibkrConfig := config.GetConfig()

	httpClient := ibkr_http.NewIBKRHttpClient()

	wsClient, err := ibkr_websocket.NewIBKRWebsocketClient(ibkrConfig.AppConfig.IBKR.Websocket.Endpoint)
	if err != nil {
		panic(err)
	}
	scheduler := gocron.NewScheduler(time.UTC)

	// Ping session to keep the session alive every minute
	scheduler.Every(1).Minutes().Do(wsClient.PingSession)
	scheduler.StartAsync()

	return &IBKRClient{
		HttpClient:      httpClient,
		WsClient:        wsClient,
		Scheduler:       scheduler,
		Config:          &config.ArbitrageConfig,
		AuthenticatedCh: make(chan bool),
	}
}
