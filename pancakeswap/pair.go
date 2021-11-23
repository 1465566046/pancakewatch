package pancakeswap

import "encoding/json"

// Pair represents a pair of two different PCS tokens.
type Pair struct {
	PairAddress  string `json:"pair_address"`
	BaseName     string `json:"base_name"`
	BaseSymbol   string `json:"base_symbol"`
	BaseAddress  string `json:"base_address"`
	QuoteName    string `json:"quote_name"`
	QuoteSymbol  string `json:"quote_symbol"`
	QuoteAddress string `json:"quote_address"`
	Price        string `json:"price"`
	BaseVolume   string `json:"base_volume"`
	QuoteVolume  string `json:"quote_volume"`
	Liquidity    string `json:"liquidity"`
	LiquidityBNB string `json:"liquidity_BNB"`
}

// SummaryRequest represents an API request of a map of pairs.
// Pair addresses in the map are separated by underscores.
type SummaryRequest struct {
	UpdatedAt int             `json:"updated_at"`
	Data      map[string]Pair `json:"data"`
}

// PairsRequest represents an API request of a map of pairs data.
// Pair addresses in the map are separated by underscores.
type PairsRequest struct {
	UpdatedAt int             `json:"updated_at"`
	Data      map[string]Pair `json:"data"`
}

// Summary returns data for the top ~1000 PancakeSwap pairs, sorted by reserves.
func Summary() (r SummaryRequest, err error) {
	data, err := apiRequest("summary")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	return
}

// Pairs returns data for the top ~1000 PancakeSwap pairs,
// sorted by reserves, with the pairs separated by underscores
func Pairs() (r PairsRequest, err error) {
	data, err := apiRequest("pairs")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	return
}
