package ibkr_websocket

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type IBKRWebsocketClient struct {
	conn          *websocket.Conn
	url           string
	done          chan struct{}
	Authenticated chan bool
}

type StsMessage struct {
	Topic string `json:"topic"`
	Args  struct {
		Authenticated bool   `json:"authenticated"`
		Competing     bool   `json:"competing"`
		Message       string `json:"message"`
		Fail          string `json:"fail"`
		ServerName    string `json:"serverName"`
		ServerVersion string `json:"serverVersion"`
		Username      string `json:"username"`
	} `json:"args"`
}

func NewIBKRWebsocketClient(wsUrl string) (*IBKRWebsocketClient, error) {
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
		done:          make(chan struct{})}, nil
}

func (client *IBKRWebsocketClient) Close() {
	client.conn.Close()
}

func (client *IBKRWebsocketClient) GetConn() *websocket.Conn {
	return client.conn
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
	err := client.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}

func (client *IBKRWebsocketClient) StartListening() {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Printf("received: %s", message)

		var stsMsg StsMessage
		if err := json.Unmarshal(message, &stsMsg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// Check authenticated flag
		if stsMsg.Topic == "sts" && stsMsg.Args.Authenticated {
			fmt.Println("Authenticated! Proceeding to subscribe to historical data.")

			client.Authenticated <- true

			// request := HistoricalDataRequest{
			// 	Conid:  "677037676",
			// 	Period: "1d",
			// 	Bar:    "1hour",
			// 	Source: "trades",
			// 	Format: "%o/%c/%h/%l",
			// }

			// client.SubscribeHistoricalData(request)
		}
	}
}
