package accounts

import "time"

const (
	TypeCash          = "cash"
	TypeBank          = "bank"
	TypeVirtualWallet = "virtual_wallet"
	TypeCreditCard    = "credit_card"
	TypeSavings       = "savings"
)

var AllowedTypes = map[string]struct{}{
	TypeCash:          {},
	TypeBank:          {},
	TypeVirtualWallet: {},
	TypeCreditCard:    {},
	TypeSavings:       {},
}

type Account struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"-"`
	Name           string     `json:"name"`
	NameNormalized string     `json:"-"`
	Type           string     `json:"type"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type CreateAccountInput struct {
	Name           string  `json:"name" binding:"required"`
	NameNormalized string  `json:"-"`
	Type           string  `json:"type" binding:"required"`
	InitialAmount  float64 `json:"initial_amount"`
	UserID         int64   `json:"-"`
	IsActive       bool    `json:"-"`
}

type UpdateAccountInput struct {
	Name           string `json:"name" binding:"required"`
	NameNormalized string `json:"-"`
	Type           string `json:"type" binding:"required"`
	IsActive       *bool  `json:"is_active,omitempty"`
}
