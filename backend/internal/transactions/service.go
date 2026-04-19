package transactions

import (
	"cenz/backend/internal/accounts"
	"cenz/backend/internal/shared/normalize"
	"errors"
	"strings"
)

var ErrAccountNotFound = errors.New("la cuenta indicada no existe")
var ErrAccountInactive = errors.New("la cuenta indicada está inactiva")
var ErrTransactionNotFound = errors.New("la transacción no existe")
var ErrCannotUpdateTransfer = errors.New("las transacciones de transferencia no se editan individualmente")
var ErrTransactionTitleRequired = errors.New("el título de la transacción es obligatorio")
var ErrTransactionCategoryRequired = errors.New("la categoría de la transacción es obligatoria")
var ErrTransactionCategoryInvalid = errors.New("la categoría de la transacción no es válida")
var ErrTransactionCategoryTypeMismatch = errors.New("la categoría no corresponde al tipo seleccionado")
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
	ListCategories() []TransactionCategoryOption
	Create(userID int64, input CreateTransactionInput) (Transaction, error)
	CreateInitialBalance(userID, accountID int64, amount float64) error
	GetAll(userID int64, filters TransactionFilters) ([]Transaction, error)
	GetByID(userID, id int64) (Transaction, error)
	Update(userID, id int64, input UpdateTransactionInput) (Transaction, error)
	Delete(userID, id int64) error
	GetBalance(userID int64) (float64, error)
	GetBalanceByAccountDetailed(userID int64) ([]AccountBalance, float64, error)
	GetSummary(userID int64, filters TransactionFilters) (TransactionSummary, error)
	GetDashboard(userID int64, filters TransactionFilters) (DashboardAnalytics, error)
	GetDashboardCategoryDetail(userID int64, filters TransactionFilters, categoryKey string) (DashboardCategoryDetail, error)
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

func (s *service) ListCategories() []TransactionCategoryOption {
	return listTransactionCategories()
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
	transaction, err := s.repo.Create(userID, prepared)
	if err != nil {
		return Transaction{}, err
	}

	return normalizeTransactionForResponse(transaction), nil
}

func (s *service) GetAll(userID int64, filters TransactionFilters) ([]Transaction, error) {
	return s.getTransactions(userID, filters)
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
	transactions, err := s.getTransactions(userID, filters)
	if err != nil {
		return TransactionSummary{}, err
	}

	summary := summarizeTransactions(transactions)
	return TransactionSummary{
		IncomeTotal:       summary.IncomeTotal,
		ExpenseTotal:      summary.ExpenseTotal,
		NetBalance:        summary.NetBalance,
		TransactionsCount: summary.TransactionsCount,
	}, nil
}

func (s *service) CreateInitialBalance(userID, accountID int64, amount float64) error {
	normalizedAmount, err := normalize.Money(amount)
	if err != nil {
		return ErrTransactionAmountInvalid
	}

	_, err = s.repo.Create(userID, CreateTransactionInput{
		UserID:      userID,
		Title:       InitialBalanceTitle,
		Amount:      normalizedAmount,
		Type:        Income,
		AccountID:   accountID,
		Category:    CategoryOther,
		Description: InitialBalanceDescription,
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
	return normalizeTransactionForResponse(transaction), nil
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

	transaction, err := s.repo.Update(userID, id, prepared)
	if err != nil {
		return Transaction{}, err
	}

	return normalizeTransactionForResponse(transaction), nil
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

	if input.Type != Income && input.Type != Expense {
		return CreateTransactionInput{}, ErrTransactionTypeInvalid
	}

	category, err := normalizeTransactionCategoryInput(input.Category, input.Type)
	if err != nil {
		return CreateTransactionInput{}, err
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

	if input.Type != Income && input.Type != Expense {
		return UpdateTransactionInput{}, ErrTransactionTypeInvalid
	}

	category, err := normalizeTransactionCategoryInput(input.Category, input.Type)
	if err != nil {
		return UpdateTransactionInput{}, err
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

func (s *service) getTransactions(userID int64, filters TransactionFilters) ([]Transaction, error) {
	baseFilters, categoryFilter := filters, filters.Category
	baseFilters.Category = nil

	var normalizedCategoryFilter *string
	if categoryFilter != nil {
		categoryKey, err := normalizeTransactionCategoryFilter(*categoryFilter)
		if err != nil {
			return nil, err
		}
		normalizedCategoryFilter = &categoryKey
	}

	transactions, err := s.repo.GetAll(userID, baseFilters)
	if err != nil {
		return nil, err
	}

	result := make([]Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		normalizedTransaction := normalizeTransactionForResponse(transaction)

		if normalizedCategoryFilter != nil && normalizedTransaction.Category != *normalizedCategoryFilter {
			continue
		}

		result = append(result, normalizedTransaction)
	}

	return result, nil
}

func normalizeTransactionForResponse(transaction Transaction) Transaction {
	transaction.Category = normalizeStoredTransactionCategory(transaction)
	transaction.CategoryLabel = resolveTransactionCategoryLabel(transaction.Category)
	return transaction
}
