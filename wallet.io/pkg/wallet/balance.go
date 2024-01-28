package wallet

type Balance struct {
	Amount map[string]int64 `json:"balance" bson:"balance"`
}
