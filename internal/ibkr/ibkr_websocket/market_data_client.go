package ibkr_websocket

import (
	"fmt"
	"strings"
)

type StreamingDataRequest struct {
	Conid  string
	Fields []string
}

type HistoricalDataRequest struct {
	Conid      string
	Exchange   string
	Period     string
	Bar        string
	OutsideRth bool
	Source     string
	Format     string
}

func (client *IBKRWebsocketClient) SubscribeStreamingData(request StreamingDataRequest) error {
	fieldsStr := `["` + strings.Join(request.Fields, `","`) + `"]`

	message := []byte(`smd+` + request.Conid + `+{
		"fields":` + fieldsStr + `
	}`)

	fmt.Println(string(message))

	err := client.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (client *IBKRWebsocketClient) SubscribeHistoricalData(request HistoricalDataRequest) error {
	message := []byte(`smh+` + request.Conid + `+{
		"period": "` + request.Period + `",
		"bar": "` + request.Bar + `", 
		"source": "` + request.Source + `", 
		"format": "` + request.Format + `"
	}`)

	err := client.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (client *IBKRWebsocketClient) SubscribeLiveOrderUpdate() error {
	message := []byte(`sor+{}`)

	err := client.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (client *IBKRWebsocketClient) SubscribeTrades() error {
	message := []byte(`str+{}`)

	err := client.Write(message)
	if err != nil {
		return err
	}

	return nil
}
