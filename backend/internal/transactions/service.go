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
	ExistsByID(userID, id int64) (bool, error)
}

type AccountBalance struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Service interface {
	Create(userID int64, input CreateTransactionInput) (Transaction, error)
	CreateInitialBalance(userID, accountID int64, amount float64) error
	GetAll(userID int64, filters TransactionFilters) ([]Transaction, error)
	GetByID(userID, id int64) (Transaction, error)
	Update(userID, id int64, input UpdateTransactionInput) (Transaction, error)
	Delete(userID, id int64) error
	GetBalance(userID int64) (float64, error)
	GetBalanceByAccountDetailed(userID int64) ([]AccountBalance, float64, error)
	GetSummary(userID int64, filters TransactionFilters) (TransactionSummary, error)
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

func (s *service) Create(userID int64, input CreateTransactionInput) (Transaction, error) {
	exists, err := s.accountFinder.ExistsByID(userID, input.AccountID)
	if err != nil {
		return Transaction{}, err
	}

	if !exists {
		return Transaction{}, ErrAccountNotFound
	}

	input.UserID = userID
	return s.repo.Create(userID, input)
}

func (s *service) GetAll(userID int64, filters TransactionFilters) ([]Transaction, error) {
	return s.repo.GetAll(userID, filters)
}

func (s *service) GetBalance(userID int64) (float64, error) {
	return s.repo.GetBalance(userID)
}

func (s *service) GetBalanceByAccountDetailed(userID int64) ([]AccountBalance, float64, error) {
	rawBalances, err := s.repo.GetBalanceByAccount(userID)
	if err != nil {
		return nil, 0, err
	}

	accountsList, err := s.accountSvc.GetAll(userID)
	if err != nil {
		return nil, 0, err
	}

	result := make([]AccountBalance, 0, len(accountsList))

	for _, acc := range accountsList {
		result = append(result, AccountBalance{
			ID:      acc.ID,
			Name:    acc.Name,
			Balance: rawBalances[acc.ID],
		})
	}

	total, err := s.repo.GetBalance(userID)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (s *service) GetSummary(userID int64, filters TransactionFilters) (TransactionSummary, error) {
	return s.repo.GetSummary(userID, filters)
}

func (s *service) CreateInitialBalance(userID, accountID int64, amount float64) error {
	_, err := s.repo.Create(userID, CreateTransactionInput{
		UserID:      userID,
		Title:       "Saldo inicial",
		Amount:      amount,
		Type:        Income,
		AccountID:   accountID,
		Category:    InitialBalanceCategory,
		Description: "Carga inicial de saldo",
	})
	return err
}

func (s *service) GetByID(userID, id int64) (Transaction, error) {
	transaction, found, err := s.repo.GetByID(userID, id)
	if err != nil {
		return Transaction{}, err
	}
	if !found {
		return Transaction{}, ErrTransactionNotFound
	}
	return transaction, nil
}

func (s *service) Update(userID, id int64, input UpdateTransactionInput) (Transaction, error) {
	exists, err := s.accountFinder.ExistsByID(userID, input.AccountID)
	if err != nil {
		return Transaction{}, err
	}
	if !exists {
		return Transaction{}, ErrAccountNotFound
	}

	current, found, err := s.repo.GetByID(userID, id)
	if err != nil {
		return Transaction{}, err
	}
	if !found {
		return Transaction{}, ErrTransactionNotFound
	}

	if current.TransferID != nil {
		return Transaction{}, ErrCannotUpdateTransfer
	}

	return s.repo.Update(userID, id, input)
}

func (s *service) Delete(userID, id int64) error {
	current, found, err := s.repo.GetByID(userID, id)
	if err != nil {
		return err
	}
	if !found {
		return ErrTransactionNotFound
	}

	if current.TransferID != nil {
		return s.repo.DeleteByTransferID(userID, *current.TransferID)
	}

	return s.repo.Delete(userID, id)
}
