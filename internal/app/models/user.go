package models

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	CoinBalance int64  `json:"coin_balance"`
}
