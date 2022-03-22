package coinbase

import (
	"context"
	"currency-vwap/services/websocket"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	//Perform Negative Testing
	_, err := NewClient("")
	require.Error(t, err)

	//Perform Positive Case
	_, err = NewClient(URL)
	require.NoError(t, err)

}

//Test to check if invalid trading pair should get error
func TestSubscribe(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		TradingPairs        []string
		ExpectedErrorResult bool
	}{
		{
			Name:                "InvalidTradingPair",
			TradingPairs:        []string{"xxx-INR"},
			ExpectedErrorResult: true,
		},
		{
			Name:                "ValidTradingPairs",
			TradingPairs:        []string{"BTC-USD"},
			ExpectedErrorResult: false,
		},
	}
	ctx := context.Background()

	response := make(chan websocket.Response)

	for _, test := range testCases {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			socketClient, err := NewClient(URL)
			require.NoError(t, err)
			err = socketClient.Subscribe(ctx, test.TradingPairs, response)
			if test.ExpectedErrorResult {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})

	}
}
