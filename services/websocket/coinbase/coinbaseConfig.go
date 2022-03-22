package coinbase

const (
	URL         = "wss://ws-feed.exchange.coinbase.com"
	RequestType = "subscribe"
	ChannelType = "matches"
)

type Channel struct {
	Name       string
	ProductIDs []string
}

// Type request will be used to send request to coinbase
type Request struct {
	Type       string    `json:"type"`
	ProductIDs []string  `json:"product_ids"`
	Channels   []Channel `json:"channels"`
}

// Type response to receive response from coinbase
type Response struct {
	Type      string    `json:"type"`
	Channels  []Channel `json:"channels"`
	Message   string    `json:"message,omitempty"`
	Size      string    `json:"size"`
	Price     string    `json:"price"`
	ProductID string    `json:"product_id"`
}
