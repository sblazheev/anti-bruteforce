package dto

type Request struct {
	Login    string `json:"login" example:"login"`
	Password string `json:"password" example:"password"`
	IP       string `json:"ip" example:"127.0.0.1"`
} // @name Event .

func New(login string, password string, ip string) (*Request, error) {
	return &Request{
		Login:    login,
		Password: password,
		IP:       ip,
	}, nil
}
