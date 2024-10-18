package ibkr_websocket

import "fmt"

func (client *IBKRWebsocketClient) PingSession() {
	message := []byte(`tic`)

	fmt.Println("Sent: ", string(message))

	client.Write(message)
}
