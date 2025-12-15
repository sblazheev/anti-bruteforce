package dto

type Event struct {
	Login    string `json:"login" example:""`
	Password string `json:"password" example:"pass"`
	IP       string `json:"ip" example:"Description"`
} // @name Event .

func New(login string, password string, ip string) (*Event, error) {
	return &Event{
		Login:    login,
		Password: password,
		IP:       ip,
	}, nil
}
