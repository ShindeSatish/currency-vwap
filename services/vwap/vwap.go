package vwap

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"
)

type DataPoint struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
}

type ProductDataMapping struct {
	SumVolumeWeighted decimal.Decimal
	SumVolume         decimal.Decimal
	VWAP              decimal.Decimal
	DataPoints        []DataPoint
}

//This is the queue that holds the data for every product/currency pair
type DataQueue struct {
	ProductInfo map[string]ProductDataMapping
	MaxSize     uint
	mu          sync.Mutex
}

func NewDataQueue(dataPoints []DataPoint, maxSize uint) (DataQueue, error) {
	if len(dataPoints) > int(maxSize) {
		return DataQueue{}, errors.New("Datapoint is exceded the maz limit")
	}
	return DataQueue{
		ProductInfo: make(map[string]ProductDataMapping),
		MaxSize:     maxSize,
		mu:          sync.Mutex{},
	}, nil

}

/**
  Function to store datapoints and calculate the VWAP value for last 200 datapoints
  We don't have to process the complete queue data every time. So, store the sum value and calulating with the sum value.
  Implemented sliding window with Queue.
*/
func (queue *DataQueue) PushData(productId string, dp DataPoint) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	//Check the productInfo available in the queue map for the requested Product Id
	if _, ok := queue.ProductInfo[productId]; ok {
		productInfo := queue.ProductInfo[productId]

		//If the datapoints count is matching with the max size
		//Remove the first data point from the list
		if uint(len(productInfo.DataPoints)) == queue.MaxSize {
			dp := productInfo.DataPoints[0]
			productInfo.DataPoints = productInfo.DataPoints[1:]
			productInfo.SumVolumeWeighted = productInfo.SumVolumeWeighted.Sub(dp.Price.Mul(dp.Volume))
			productInfo.SumVolume = productInfo.SumVolume.Sub(dp.Volume)
			if !productInfo.SumVolume.IsZero() {
				//Calculate VWAP
				productInfo.VWAP = productInfo.SumVolumeWeighted.Div(productInfo.SumVolume)
			}

			queue.ProductInfo[productId] = productInfo
		}

		//Calculate the VWAP and append the new datapoint into the ProductInfo datapoints queue
		productInfo.SumVolumeWeighted = productInfo.SumVolumeWeighted.Add(dp.Price.Mul(dp.Volume))
		productInfo.SumVolume = productInfo.SumVolume.Add(dp.Volume)
		productInfo.VWAP = productInfo.SumVolumeWeighted.Div(productInfo.SumVolume)
		productInfo.DataPoints = append(productInfo.DataPoints, dp)
		queue.ProductInfo[productId] = productInfo

	} else {

		//Calculate and store the initial VWAP
		valueWeighted := dp.Price.Mul(dp.Volume)
		productInfo := ProductDataMapping{
			SumVolumeWeighted: valueWeighted,
			SumVolume:         dp.Volume,
			VWAP:              valueWeighted.Div(dp.Volume),
		}
		productInfo.DataPoints = append(productInfo.DataPoints, dp)
		queue.ProductInfo[productId] = productInfo
	}
}
