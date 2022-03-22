# Currency VWAP Calculator
Calculate VWAP from crypto currency data feed.
Currently Used coinbase as realtime data provider.

# Getting Started
First, make sure that you have go installed on your machine. Then ...
```
make run
```
or 
```
make build
``` 
then ... run the application with 
```
./currency-vwap -url "<coinbase_websocket_url>" -tradingPairs "<crypto cyrrency pairs>" -windowSize <Window size>
```

## Config
The following flags are available while running the project through CLI
* `tradingPairs`: a comma separated strings of crypto currencies pairs, default value is  `BTC-USD,ETH-USD,ETH-BTC`
* `url`: the URL of the websocket server to use, default is coinbase URL.
* `windowSize`: the sliding window used for the VWAP calculation, default value is 200.

## Design
 main.go is the entry point.
 
 `./services/service.go` is the service. and this service used two components.
 * A websocket client that pulls data off an exchange.
    * The default choice is coinbase.
    * Any exchange can be used as long as it implements the client interface defined in the websocket package.
 * Datapoint queue defined in in VWAP package
    * Every time the data pushed into the queue VWAP will be calculated with currency pair
    * VWAP Calculated in constant time O(1).
    * Storing the calculated values so no loop over datapoints required.

## Tests
* This command will run all the test cases.
```
make test
```

# Result 
* For printing VWAP used service.PrintVWAP function. we can change the print formatting here.

# Packages Used 
* For precision I used https://github.com/shopspring/decimal for all calculation.
* https://github.com/stretchr/testify Used for writing unit tests
* https://golang.org/x/xerrors For better error formating 
