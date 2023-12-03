package cubawheeler

type LoginRequest struct {
	Email string  `json:"email"`
	Otp   *string `json:"otp,omitempty"`
	Pin   *string `json:"pin,omitempty"`
}
