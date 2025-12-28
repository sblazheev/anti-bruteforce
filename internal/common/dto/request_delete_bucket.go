package dto

type RequestDeleteBucket struct {
	Login string `json:"login" example:"login"`
	IP    string `json:"ip" example:"127.0.0.1"`
} // @name RequestDeleteBucket .

func NewDeleteBucket(login string, ip string) (*RequestDeleteBucket, error) {
	return &RequestDeleteBucket{
		Login: login,
		IP:    ip,
	}, nil
}
