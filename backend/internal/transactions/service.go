package transactions

import (
	"cartera-app/backend/internal/accounts"
	"errors"
)

var ErrAccountNotFound = errors.New("la cuenta indicada no existe")

type AccountFinder interface {
	ExistsByID(id int64) bool
}

type AccountBalance struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Service interface {
	Create(input CreateTransactionInput) (Transaction, error)
	GetAll() []Transaction
	GetBalance() float64
	GetBalanceByAccountDetailed() ([]AccountBalance, float64)
}

type service struct {
	repo          Repository
	accountFinder AccountFinder
	accountSvc    accounts.Service
}

func NewService(repo Repository, accountFinder AccountFinder, accountSvc accounts.Service) Service {
	return &service{
		repo:          repo,
		accountFinder: accountFinder,
		accountSvc:    accountSvc,
	}
}

func (s *service) Create(input CreateTransactionInput) (Transaction, error) {
	if !s.accountFinder.ExistsByID(input.AccountID) {
		return Transaction{}, ErrAccountNotFound
	}

	transaction := s.repo.Create(input)
	return transaction, nil
}

func (s *service) GetAll() []Transaction {
	return s.repo.GetAll()
}

func (s *service) GetBalance() float64 {
	return s.repo.GetBalance()
}

func (s *service) GetBalanceByAccountDetailed() ([]AccountBalance, float64) {
	rawBalances := s.repo.GetBalanceByAccount()
	accountsList := s.accountSvc.GetAll()

	var result []AccountBalance

	for _, acc := range accountsList {
		balance := rawBalances[acc.ID]

		result = append(result, AccountBalance{
			ID:      acc.ID,
			Name:    acc.Name,
			Balance: balance,
		})
	}

	return result, s.repo.GetBalance()
}
