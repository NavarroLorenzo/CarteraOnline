package transactions

import (
	"cartera-app/backend/internal/accounts"
	"cartera-app/backend/internal/shared/normalize"
	"errors"
	"strings"
)

var ErrAccountNotFound = errors.New("la cuenta indicada no existe")
var ErrAccountInactive = errors.New("la cuenta indicada está inactiva")
var ErrTransactionNotFound = errors.New("la transacción no existe")
var ErrCannotUpdateTransfer = errors.New("las transacciones de transferencia no se editan individualmente")
var ErrTransactionTitleRequired = errors.New("el título de la transacción es obligatorio")
var ErrTransactionCategoryRequired = errors.New("la categoría de la transacción es obligatoria")
var ErrTransactionAmountInvalid = errors.New("el monto debe ser mayor a cero")
var ErrTransactionTypeInvalid = errors.New("el tipo de transacción no es válido")

type AccountFinder interface {
	GetByID(userID, id int64) (accounts.Account, bool, error)
	GetAll(userID int64) ([]accounts.Account, error)
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
	prepared, err := normalizeCreateInput(input)
	if err != nil {
		return Transaction{}, err
	}

	account, found, err := s.accountFinder.GetByID(userID, prepared.AccountID)
	if err != nil {
		return Transaction{}, err
	}

	if !found {
		return Transaction{}, ErrAccountNotFound
	}
	if !account.IsActive {
		return Transaction{}, ErrAccountInactive
	}

	prepared.UserID = userID
	return s.repo.Create(userID, prepared)
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
	normalizedAmount, err := normalize.Money(amount)
	if err != nil {
		return ErrTransactionAmountInvalid
	}

	_, err = s.repo.Create(userID, CreateTransactionInput{
		UserID:      userID,
		Title:       "Saldo inicial",
		Amount:      normalizedAmount,
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

	prepared, err := normalizeUpdateInput(input)
	if err != nil {
		return Transaction{}, err
	}

	account, found, err := s.accountFinder.GetByID(userID, prepared.AccountID)
	if err != nil {
		return Transaction{}, err
	}
	if !found {
		return Transaction{}, ErrAccountNotFound
	}
	if !account.IsActive {
		return Transaction{}, ErrAccountInactive
	}

	return s.repo.Update(userID, id, prepared)
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

func normalizeCreateInput(input CreateTransactionInput) (CreateTransactionInput, error) {
	title := normalize.Optional(input.Title)
	if title == "" {
		return CreateTransactionInput{}, ErrTransactionTitleRequired
	}

	category := normalize.Optional(input.Category)
	if category == "" {
		return CreateTransactionInput{}, ErrTransactionCategoryRequired
	}

	if input.Type != Income && input.Type != Expense {
		return CreateTransactionInput{}, ErrTransactionTypeInvalid
	}

	amount, err := normalize.Money(input.Amount)
	if err != nil {
		return CreateTransactionInput{}, ErrTransactionAmountInvalid
	}

	return CreateTransactionInput{
		UserID:      input.UserID,
		Title:       title,
		Amount:      amount,
		Type:        input.Type,
		AccountID:   input.AccountID,
		Category:    category,
		Description: normalize.Optional(input.Description),
		TransferID:  input.TransferID,
	}, nil
}

func normalizeUpdateInput(input UpdateTransactionInput) (UpdateTransactionInput, error) {
	title := normalize.Optional(input.Title)
	if title == "" {
		return UpdateTransactionInput{}, ErrTransactionTitleRequired
	}

	category := normalize.Optional(input.Category)
	if category == "" {
		return UpdateTransactionInput{}, ErrTransactionCategoryRequired
	}

	if input.Type != Income && input.Type != Expense {
		return UpdateTransactionInput{}, ErrTransactionTypeInvalid
	}

	amount, err := normalize.Money(input.Amount)
	if err != nil {
		return UpdateTransactionInput{}, ErrTransactionAmountInvalid
	}

	return UpdateTransactionInput{
		Title:       title,
		Amount:      amount,
		Type:        input.Type,
		AccountID:   input.AccountID,
		Category:    category,
		Description: strings.TrimSpace(normalize.Optional(input.Description)),
	}, nil
}
