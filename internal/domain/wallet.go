package domain

type Wallet struct {
	Address      string        `json:"address"`
	SecretKey    string        `json:"secret-key,omitempty"`
	Balance      float64       `json:"balance"`
	Transactions []Transaction `json:"transactions,omitempty"`
}
