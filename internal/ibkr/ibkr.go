package ibkr

import (
	"time"

	"github.com/btc-etf-arbitrage/internal/config"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_http"
	"github.com/btc-etf-arbitrage/internal/ibkr/ibkr_websocket"
	"github.com/rs/zerolog"

	"github.com/go-co-op/gocron"
)

type IbkrClient struct {
	HttpClient      *ibkr_http.IBKRHttpClient
	WsClient        *ibkr_websocket.IBKRWebsocketClient
	Scheduler       *gocron.Scheduler
	Config          *config.AppConfig
	logger          zerolog.Logger
	AuthenticatedCh chan bool
}

// NewIBKRClient creates a new instance of IbkrClient
func NewIBKRClient(logger zerolog.Logger) *IbkrClient {
	appConfig := config.GetAppConfig()
	httpClient := ibkr_http.NewIBKRHttpClient()
	wsClient, err := ibkr_websocket.NewIBKRWebsocketClient(appConfig.IbkrConfig.WsEndpoint)
	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	// Ping session to keep the session alive every minute
	scheduler.Every(1).Minutes().Do(wsClient.PingSession)
	scheduler.StartAsync()

	return &IbkrClient{
		HttpClient:      httpClient,
		WsClient:        wsClient,
		Scheduler:       scheduler,
		Config:          &appConfig,
		AuthenticatedCh: make(chan bool),
		logger:          logger.With().Str("component", "ibkr").Logger(),
	}
}
