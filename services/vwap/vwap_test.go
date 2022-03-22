package vwap

import (
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestMaxWindowSize(t *testing.T) {
	t.Parallel()
	queue, err := NewDataQueue([]DataPoint{}, 5)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(10)
	Product := "BTC-USD"
	for counter := 1; counter <= 10; counter++ {
		go func() {
			data := DataPoint{Price: decimal.NewFromInt(10), Volume: decimal.NewFromInt(30)}
			queue.PushData(Product, data)
			wg.Done()
		}()
	}

	wg.Wait()
	require.Len(t, queue.ProductInfo[Product].DataPoints, 5)
}

//Test To push the data into the queue and calcualte VWAP. And validate the calculated VWAP is correct
func TestVWAP(t *testing.T) {
	t.Parallel()

	//Define Test Cases
	testCases := []struct {
		Name         string
		DataPoint    []DataPoint
		ProductId    string
		ExpectedVWAP decimal.Decimal
		MaxSize      uint
	}{
		{
			Name:         "EmptyDataPoints",
			DataPoint:    []DataPoint{},
			ProductId:    "BTC-USD",
			ExpectedVWAP: decimal.Decimal{},
		},
		{
			Name: "VALIDATE-BTC-USD",
			DataPoint: []DataPoint{
				{Price: decimal.NewFromInt(11), Volume: decimal.NewFromInt(30)},
				{Price: decimal.NewFromInt(12), Volume: decimal.NewFromInt(60)},
				{Price: decimal.NewFromInt(10), Volume: decimal.NewFromInt(40)},
				{Price: decimal.NewFromInt(15), Volume: decimal.NewFromInt(50)},
			},
			ProductId:    "BTC-USD",
			ExpectedVWAP: decimal.RequireFromString("12.2222222222222222"),
			MaxSize:      4,
		},
		{
			Name: "VALIDATE-ETH-USD",
			DataPoint: []DataPoint{
				{Price: decimal.NewFromInt(21), Volume: decimal.NewFromInt(10)},
				{Price: decimal.NewFromInt(23), Volume: decimal.NewFromInt(30)},
				{Price: decimal.NewFromInt(22), Volume: decimal.NewFromInt(20)},
				{Price: decimal.NewFromInt(25), Volume: decimal.NewFromInt(60)},
			},
			ProductId:    "ETH_USD",
			ExpectedVWAP: decimal.RequireFromString("23.6666666666666667"),
			MaxSize:      4,
		},
		{
			Name: "VALIDATE-ETH-BTC",
			DataPoint: []DataPoint{
				{Price: decimal.NewFromInt(20), Volume: decimal.RequireFromString("11")},
				{Price: decimal.NewFromInt(20), Volume: decimal.RequireFromString("50")},
			},
			ProductId:    "ETH-BTC",
			ExpectedVWAP: decimal.RequireFromString("20"),
			MaxSize:      4,
		},
	}

	//Evaluate Test Cases
	for _, test := range testCases {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			dQueue, err := NewDataQueue([]DataPoint{}, test.MaxSize)
			require.NoError(t, err)

			for _, d := range test.DataPoint {
				dQueue.PushData(test.ProductId, d)
			}
			require.Equal(t, test.ExpectedVWAP.String(), dQueue.ProductInfo[test.ProductId].VWAP.String())
		})
	}
}
