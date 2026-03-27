package accounts

import "time"

type Account struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountInput struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}
