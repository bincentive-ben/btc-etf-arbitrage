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

func (client *IBKRWebsocketClient) SubscribeStreamingData(request StreamingDataRequest) {
	fieldsStr := `["` + strings.Join(request.Fields, `","`) + `"]`

	message := []byte(`smd+` + request.Conid + `+{
		"fields":` + fieldsStr + `
	}`)

	fmt.Println(string(message))

	client.Write(message)
}

func (client *IBKRWebsocketClient) SubscribeHistoricalData(request HistoricalDataRequest) {
	message := []byte(`smh+` + request.Conid + `+{
		"period": "` + request.Period + `",
		"bar": "` + request.Bar + `", 
		"source": "` + request.Source + `", 
		"format": "` + request.Format + `"
	}`)

	client.Write(message)
}
