package coinbase

import (
	"context"
	"currency-vwap/services/websocket"
	"encoding/json"
	"errors"
	"log"

	wsocket "golang.org/x/net/websocket"
	"golang.org/x/xerrors"
)

type client struct {
	conn *wsocket.Conn
}

// wesocketResponse converts the coinbase response into a websocket response.
func socketResponse(res Response) websocket.Response {
	return websocket.Response{
		Type:      res.Type,
		Size:      res.Size,
		Price:     res.Price,
		ProductID: res.ProductID,
	}
}

// Subscribe subscribes to the websocket.
func (c *client) Subscribe(ctx context.Context, tradingPairs []string, receiver chan websocket.Response) error {
	if len(tradingPairs) == 0 {
		return errors.New("Trading Pairs can not be empty")
	}

	subscription := Request{
		Type:       RequestType,
		ProductIDs: tradingPairs,
		Channels: []Channel{
			{Name: ChannelType},
		},
	}

	requestData, err := json.Marshal(subscription)
	if err != nil {
		return xerrors.Errorf("Failed to marshal: %w", err)
	}

	err = wsocket.Message.Send(c.conn, requestData)
	if err != nil {
		return xerrors.Errorf("Failed at sending subscription: %w", err)
	}

	var response Response

	err = wsocket.JSON.Receive(c.conn, &response)
	if err != nil {
		return xerrors.Errorf("Failed to get response: %w", err)
	}

	if response.Type == "error" {
		return xerrors.Errorf("Response type error with subscription: %w", response.Message)
	}
	/**
		keep listning to the websocket data
	**/
	go func() {
		for {
			var message Response
			err := wsocket.JSON.Receive(c.conn, &message)
			if err != nil {
				log.Printf("Failed receiving message: %s", err)
				break
			}
			receiver <- socketResponse(message)
		}
	}()

	return nil
}

// NewClient returns a new websocket client.
func NewClient(url string) (websocket.Client, error) {
	conn, err := wsocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	log.Printf("websocket connected to: %s", url)

	return &client{
		conn: conn,
	}, nil
}
