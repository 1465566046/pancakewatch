package pancakeswap

import (
	"encoding/json"
)

// Token represents a PCS (BSC) token.
type Token struct {
	// Name is the
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Price    string `json:"price"`
	PriceBNB string `json:"price_BNB"`
}

// TokensRequest represents an API request of a map of tokens.
type TokensRequest struct {
	UpdatedAt int              `json:"updated_at"`
	Data      map[string]Token `json:"data"`
}

// TokenRequest represents an API request of a token.
type TokenRequest struct {
	UpdatedAt int   `json:"updated_at"`
	Data      Token `json:"data"`
}

// Tokens returns the tokens in the top ~1000 pairs on PancakeSwap, sorted by reserves.
func Tokens() (r TokensRequest, err error) {
	data, err := apiRequest("tokens")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	return
}

// TokenInfo returns the token information, based on address.
func TokenInfo(address string) (r TokenRequest, err error) {
	data, err := apiRequest("tokens/" + address)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	return
}
