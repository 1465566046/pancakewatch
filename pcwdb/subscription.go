package pcwdb

// Subscription represents the information associated
// with a token address.
type Subscription struct {
	PhoneNumber   int     `json:"phoneNumber"`
	TargetPrice   float32 `json:"targetPrice"`
	IsTargetUnder bool    `json:"isTargetUnder"`
}
