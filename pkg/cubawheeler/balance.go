package cubawheeler

type Balance struct {
	Amount map[string]int64 `json:"balance" bson:"balance"`
}
