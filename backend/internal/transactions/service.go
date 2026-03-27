package transactions

import (
	"cartera-app/backend/internal/accounts"
	"errors"
)

var ErrAccountNotFound = errors.New("la cuenta indicada no existe")

type AccountFinder interface {
	ExistsByID(id int64) (bool, error)
}

type AccountBalance struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Service interface {
	Create(input CreateTransactionInput) (Transaction, error)
	GetAll(filters TransactionFilters) ([]Transaction, error)
	GetBalance() (float64, error)
	GetBalanceByAccountDetailed() ([]AccountBalance, float64, error)
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
	exists, err := s.accountFinder.ExistsByID(input.AccountID)
	if err != nil {
		return Transaction{}, err
	}

	if !exists {
		return Transaction{}, ErrAccountNotFound
	}

	return s.repo.Create(input)
}

func (s *service) GetAll(filters TransactionFilters) ([]Transaction, error) {
	return s.repo.GetAll(filters)
}

func (s *service) GetBalance() (float64, error) {
	return s.repo.GetBalance()
}

func (s *service) GetBalanceByAccountDetailed() ([]AccountBalance, float64, error) {
	rawBalances, err := s.repo.GetBalanceByAccount()
	if err != nil {
		return nil, 0, err
	}

	accountsList, err := s.accountSvc.GetAll()
	if err != nil {
		return nil, 0, err
	}

	var result []AccountBalance

	for _, acc := range accountsList {
		result = append(result, AccountBalance{
			ID:      acc.ID,
			Name:    acc.Name,
			Balance: rawBalances[acc.ID],
		})
	}

	total, err := s.repo.GetBalance()
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}
