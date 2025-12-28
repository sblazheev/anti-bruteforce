package dto

type RequestNet struct {
	Net string `json:"net" example:"127.0.0.0/24"`
} // @name RequestIP .

func NewIP(net string) (*RequestNet, error) {
	return &RequestNet{
		Net: net,
	}, nil
}
