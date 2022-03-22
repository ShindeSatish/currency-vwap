package services

import (
	"context"
	"currency-vwap/services/vwap"
	"currency-vwap/services/websocket"
	"fmt"

	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

type service struct {
	websocketClient websocket.Client
	tradingPairs    []string
	dataQueue       *vwap.DataQueue
}

func NewService(wc websocket.Client, tradingPairs []string, dataQueue *vwap.DataQueue) *service {
	return &service{
		websocketClient: wc,
		tradingPairs:    tradingPairs,
		dataQueue:       dataQueue,
	}
}

func (s *service) Run(ctx context.Context) error {
	response := make(chan websocket.Response)

	err := s.websocketClient.Subscribe(ctx, s.tradingPairs, response)
	if err != nil {
		return xerrors.Errorf("Service subscription err: %w", err)
	}

	/**
	* Range over the response channel and get collect the match data from websocket
	* Once we get the data calculate the VWAP on that data.
	**/

	for responseData := range response {
		if responseData.Price == "" {
			continue
		}

		decimalPrice, err := decimal.NewFromString(responseData.Price)
		if err != nil {
			return xerrors.Errorf("decimalPrice %s: %w", responseData.Price, err)
		}

		decimalSize, err := decimal.NewFromString(responseData.Size)
		if err != nil {
			return xerrors.Errorf("decimalSize %s: %w", responseData.Size, err)
		}

		s.dataQueue.PushData(responseData.ProductID, vwap.DataPoint{
			Price:  decimalPrice,
			Volume: decimalSize,
		})
		s.PrintVWAP(responseData.ProductID)

	}
	return nil
}

func (s *service) PrintVWAP(productId string) {
	fmt.Printf("\nVWAP for Trading Pair %s = %v", productId, s.dataQueue.ProductInfo[productId].VWAP)
}
