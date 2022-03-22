package main

import (
	"context"
	"currency-vwap/services"
	"currency-vwap/services/vwap"
	"currency-vwap/services/websocket/coinbase"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//Defined the constants
const (
	tradingPairs = "BTC-USD,ETH-USD,ETH-BTC" //Default Trading Pairs
	windowSize   = 200
)

func main() {
	ctx := context.Background()

	tradingPairs := flag.String("tradingPairs", tradingPairs, "Trading Pairs")
	URL := flag.String("url", coinbase.URL, "Web Socket URL")
	windowSize := flag.Uint("windowSize", windowSize, "window size")

	// Intercepting Exit program signals.
	go func() {
		exitProgram := make(chan os.Signal, 1)
		signal.Notify(exitProgram, syscall.SIGTERM, syscall.SIGINT)

		signal := <-exitProgram

		log.Printf("Terminating the Program %s", signal)

		os.Exit(0)
	}()

	websocketClient, err := coinbase.NewClient(*URL)
	if err != nil {
		log.Fatal(err)
	}

	dataQueue, err := vwap.NewDataQueue([]vwap.DataPoint{}, *windowSize)
	if err != nil {
		log.Fatal(err)
	}
	//Create a service and run the service
	pairs := strings.Split(*tradingPairs, ",")
	service := services.NewService(websocketClient, pairs, &dataQueue)

	err = service.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
