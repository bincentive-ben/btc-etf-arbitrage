package ibkr_websocket

func (client *IBKRWebsocketClient) PingSession() {
	message := []byte(`tic`)

	client.Write(message)
}
