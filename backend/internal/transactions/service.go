package transactions

import (
	"cartera-app/backend/internal/accounts"
	"errors"
)

var ErrAccountNotFound = errors.New("la cuenta indicada no existe")
var ErrTransactionNotFound = errors.New("la transacción no existe")
var ErrCannotUpdateTransfer = errors.New("las transacciones de transferencia no se editan individualmente")

const InitialBalanceCategory = "initial_balance"

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
	CreateInitialBalance(accountID int64, amount float64) error
	GetAll(filters TransactionFilters) ([]Transaction, error)
	GetByID(id int64) (Transaction, error)
	Update(id int64, input UpdateTransactionInput) (Transaction, error)
	Delete(id int64) error
	GetBalance() (float64, error)
	GetBalanceByAccountDetailed() ([]AccountBalance, float64, error)
	GetSummary(filters TransactionFilters) (TransactionSummary, error)
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

func (s *service) GetSummary(filters TransactionFilters) (TransactionSummary, error) {
	return s.repo.GetSummary(filters)
}

func (s *service) CreateInitialBalance(accountID int64, amount float64) error {
	_, err := s.repo.Create(CreateTransactionInput{
		Title:       "Saldo inicial",
		Amount:      amount,
		Type:        Income,
		AccountID:   accountID,
		Category:    InitialBalanceCategory,
		Description: "Carga inicial de saldo",
	})
	return err
}

func (s *service) GetByID(id int64) (Transaction, error) {
	transaction, found, err := s.repo.GetByID(id)
	if err != nil {
		return Transaction{}, err
	}
	if !found {
		return Transaction{}, ErrTransactionNotFound
	}
	return transaction, nil
}

func (s *service) Update(id int64, input UpdateTransactionInput) (Transaction, error) {
	exists, err := s.accountFinder.ExistsByID(input.AccountID)
	if err != nil {
		return Transaction{}, err
	}
	if !exists {
		return Transaction{}, ErrAccountNotFound
	}

	current, found, err := s.repo.GetByID(id)
	if err != nil {
		return Transaction{}, err
	}
	if !found {
		return Transaction{}, ErrTransactionNotFound
	}

	if current.TransferID != nil {
		return Transaction{}, ErrCannotUpdateTransfer
	}

	return s.repo.Update(id, input)
}

func (s *service) Delete(id int64) error {
	current, found, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if !found {
		return ErrTransactionNotFound
	}

	if current.TransferID != nil {
		return s.repo.DeleteByTransferID(*current.TransferID)
	}

	return s.repo.Delete(id)
}
