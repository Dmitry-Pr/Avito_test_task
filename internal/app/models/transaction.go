package models

import "time"

type Transaction struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Type       string    `json:"type"`
	Amount     int64     `json:"amount"`
	FromUserID int64     `json:"from_user_id,omitempty"`
	ToUserID   int64     `json:"to_user_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
