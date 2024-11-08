package ibkr_websocket

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type IBKRWebsocketClient struct {
	conn          *websocket.Conn
	url           string
	done          chan struct{}
	Authenticated chan bool
	logger        zerolog.Logger
}

func NewIBKRWebsocketClient(wsUrl string, logger zerolog.Logger) (*IBKRWebsocketClient, error) {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Connect to the WebSocket server
	conn, _, err := dialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, err
	}

	return &IBKRWebsocketClient{
		conn:          conn,
		url:           wsUrl,
		Authenticated: make(chan bool),
		logger:        logger.With().Str("component", "IBKRWebsocketClient").Logger(),
		done:          make(chan struct{})}, nil
}

func (client *IBKRWebsocketClient) Reconnect() error {
	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *IBKRWebsocketClient) Close() {
	client.conn.Close()
}

func (client *IBKRWebsocketClient) GetConn() (*websocket.Conn, error) {
	if client.conn != nil {
		return client.conn, nil
	} else {
		return nil, fmt.Errorf("websocket connection is nil")
	}
}
func (client *IBKRWebsocketClient) Listen() {
	defer close(client.done)
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("received: %s", message)
	}
}

func (client *IBKRWebsocketClient) Write(message []byte) error {
	if client.conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	err := client.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}
