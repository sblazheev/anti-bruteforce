package dto

type RequestCheck struct {
	Login    string `json:"login" example:"login"`
	Password string `json:"password" example:"password"`
	IP       string `json:"ip" example:"127.0.0.1"`
} // @name RequestCheck .

func NewCheck(login string, password string, ip string) (*RequestCheck, error) {
	return &RequestCheck{
		Login:    login,
		Password: password,
		IP:       ip,
	}, nil
}
