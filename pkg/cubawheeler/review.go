package cubawheeler

type Review struct {
	ID      string  `json:"id"`
	From    string  `json:"from"`
	To      string  `json:"to"`
	Comment *string `json:"comment,omitempty"`
	Rate    float64 `json:"rate"`
}
