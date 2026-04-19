package transactions

import (
	"testing"
	"time"

	"cenz/backend/internal/accounts"
)

type fakeTransactionsRepo struct {
	transaction          Transaction
	found                bool
	createCalled         bool
	deleteCalled         bool
	deleteTransferCalled bool
	deletedTransferID    string
	lastUserIDForGetByID int64
	lastUserIDForDelete  int64
	transactions         []Transaction
}

func (f *fakeTransactionsRepo) Create(userID int64, input CreateTransactionInput) (Transaction, error) {
	f.createCalled = true
	return Transaction{
		ID:          1,
		UserID:      userID,
		Title:       input.Title,
		Amount:      input.Amount,
		Type:        input.Type,
		AccountID:   input.AccountID,
		Category:    input.Category,
		Description: input.Description,
		TransferID:  input.TransferID,
		CreatedAt:   time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC),
	}, nil
}

func (f *fakeTransactionsRepo) GetAll(userID int64, filters TransactionFilters) ([]Transaction, error) {
	result := make([]Transaction, 0)

	for _, transaction := range f.transactions {
		if transaction.UserID != 0 && transaction.UserID != userID {
			continue
		}

		if filters.AccountID != nil && transaction.AccountID != *filters.AccountID {
			continue
		}

		if filters.Type != nil && transaction.Type != *filters.Type {
			continue
		}

		if filters.Category != nil && transaction.Category != *filters.Category {
			continue
		}

		if filters.DateFrom != nil && transaction.CreatedAt.Before(*filters.DateFrom) {
			continue
		}

		if filters.DateTo != nil && transaction.CreatedAt.After(*filters.DateTo) {
			continue
		}

		result = append(result, transaction)
	}

	return result, nil
}

func (f *fakeTransactionsRepo) GetByID(userID, id int64) (Transaction, bool, error) {
	f.lastUserIDForGetByID = userID
	if userID != 1 {
		return Transaction{}, false, nil
	}

	return f.transaction, f.found, nil
}

func (f *fakeTransactionsRepo) Update(userID, id int64, input UpdateTransactionInput) (Transaction, error) {
	return Transaction{}, nil
}

func (f *fakeTransactionsRepo) Delete(userID, id int64) error {
	f.deleteCalled = true
	f.lastUserIDForDelete = userID
	return nil
}

func (f *fakeTransactionsRepo) DeleteByTransferID(userID int64, transferID string) error {
	f.deleteTransferCalled = true
	f.deletedTransferID = transferID
	f.lastUserIDForDelete = userID
	return nil
}

func (f *fakeTransactionsRepo) GetBalance(userID int64) (float64, error) {
	return 0, nil
}

func (f *fakeTransactionsRepo) GetBalanceByAccount(userID int64) (map[int64]float64, error) {
	return map[int64]float64{}, nil
}

func (f *fakeTransactionsRepo) GetSummary(userID int64, filters TransactionFilters) (TransactionSummary, error) {
	return TransactionSummary{}, nil
}

type fakeAccountsService struct {
	account  accounts.Account
	found    bool
	accounts []accounts.Account
}

func (f *fakeAccountsService) Create(userID int64, input accounts.CreateAccountInput) (accounts.Account, error) {
	return accounts.Account{}, nil
}

func (f *fakeAccountsService) GetAll(userID int64) ([]accounts.Account, error) {
	return f.accounts, nil
}

func (f *fakeAccountsService) GetByID(userID, id int64) (accounts.Account, bool, error) {
	return f.account, f.found, nil
}

func (f *fakeAccountsService) ExistsByID(userID, id int64) (bool, error) {
	return f.found, nil
}

func (f *fakeAccountsService) Update(userID, id int64, input accounts.UpdateAccountInput) (accounts.Account, error) {
	return accounts.Account{}, nil
}

func (f *fakeAccountsService) Delete(userID, id int64) error {
	return nil
}

func TestCreateRejectsInactiveAccount(t *testing.T) {
	repo := &fakeTransactionsRepo{}
	accountSvc := &fakeAccountsService{
		account: accounts.Account{ID: 5, Name: "Caja", IsActive: false},
		found:   true,
	}

	service := NewService(repo, accountSvc, accountSvc)

	_, err := service.Create(1, CreateTransactionInput{
		Title:     "Compra",
		Amount:    100,
		Type:      Expense,
		AccountID: 5,
		Category:  CategoryOther,
	})
	if err != ErrAccountInactive {
		t.Fatalf("expected ErrAccountInactive, got %v", err)
	}
	if repo.createCalled {
		t.Fatal("expected repo.Create not to be called")
	}
}

func TestDeleteTransferDeletesWholeGroup(t *testing.T) {
	transferID := "tr_test"
	repo := &fakeTransactionsRepo{
		transaction: Transaction{ID: 7, TransferID: &transferID},
		found:       true,
	}
	accountSvc := &fakeAccountsService{}

	service := NewService(repo, accountSvc, accountSvc)

	if err := service.Delete(1, 7); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.deleteTransferCalled {
		t.Fatal("expected DeleteByTransferID to be called")
	}
	if repo.deletedTransferID != transferID {
		t.Fatalf("expected transfer id %q, got %q", transferID, repo.deletedTransferID)
	}
	if repo.deleteCalled {
		t.Fatal("expected individual Delete not to be called")
	}
}

func TestDeleteOtherUserTransactionReturnsNotFound(t *testing.T) {
	repo := &fakeTransactionsRepo{
		transaction: Transaction{ID: 9, Title: "Ajena"},
		found:       true,
	}
	accountSvc := &fakeAccountsService{}

	service := NewService(repo, accountSvc, accountSvc)

	err := service.Delete(2, 9)
	if err != ErrTransactionNotFound {
		t.Fatalf("expected ErrTransactionNotFound, got %v", err)
	}
}

func TestGetDashboardBuildsAnalyticsExcludingTransfersButKeepingInitialBalance(t *testing.T) {
	repo := &fakeTransactionsRepo{
		transactions: []Transaction{
			{
				ID:        1,
				UserID:    1,
				Title:     "Sueldo",
				Amount:    3000,
				Type:      Income,
				AccountID: 10,
				Category:  "salario",
				CreatedAt: time.Date(2026, 4, 7, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				UserID:    1,
				Title:     "Supermercado",
				Amount:    500,
				Type:      Expense,
				AccountID: 10,
				Category:  "supermercado",
				CreatedAt: time.Date(2026, 4, 7, 15, 0, 0, 0, time.UTC),
			},
			{
				ID:        3,
				UserID:    1,
				Title:     "Taxi",
				Amount:    120,
				Type:      Expense,
				AccountID: 10,
				Category:  "taxi",
				CreatedAt: time.Date(2026, 4, 5, 9, 0, 0, 0, time.UTC),
			},
			{
				ID:         4,
				UserID:     1,
				Title:      "Transferencia enviada",
				Amount:     700,
				Type:       Expense,
				AccountID:  10,
				Category:   CategoryTransfer,
				TransferID: pointerToString("tr_1"),
				CreatedAt:  time.Date(2026, 4, 6, 9, 0, 0, 0, time.UTC),
			},
			{
				ID:          5,
				UserID:      1,
				Title:       InitialBalanceTitle,
				Amount:      900,
				Type:        Income,
				AccountID:   10,
				Category:    CategoryOther,
				Description: InitialBalanceDescription,
				CreatedAt:   time.Date(2026, 4, 1, 8, 0, 0, 0, time.UTC),
			},
			{
				ID:        6,
				UserID:    1,
				Title:     "Sueldo marzo",
				Amount:    2500,
				Type:      Income,
				AccountID: 10,
				Category:  "salario",
				CreatedAt: time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:        7,
				UserID:    1,
				Title:     "Internet",
				Amount:    100,
				Type:      Expense,
				AccountID: 10,
				Category:  "internet",
				CreatedAt: time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:        8,
				UserID:    1,
				Title:     "Café",
				Amount:    80,
				Type:      Expense,
				AccountID: 10,
				Category:  "comida",
				CreatedAt: time.Date(2026, 4, 10, 9, 30, 0, 0, time.UTC),
			},
		},
	}
	accountSvc := &fakeAccountsService{
		accounts: []accounts.Account{
			{ID: 10, Name: "Cuenta sueldo"},
		},
	}

	service := NewService(repo, accountSvc, accountSvc)

	dateTo := time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)
	period := DashboardPeriodMonth

	dashboard, err := service.GetDashboard(1, TransactionFilters{
		Period: &period,
		DateTo: &dateTo,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if dashboard.PeriodSummary.IncomeTotal != 3900 {
		t.Fatalf("expected income total 3900, got %v", dashboard.PeriodSummary.IncomeTotal)
	}
	if dashboard.PeriodSummary.ExpenseTotal != 620 {
		t.Fatalf("expected expense total 620, got %v", dashboard.PeriodSummary.ExpenseTotal)
	}
	if dashboard.PeriodSummary.NetBalance != 3280 {
		t.Fatalf("expected net balance 3280, got %v", dashboard.PeriodSummary.NetBalance)
	}
	if len(dashboard.ExpenseCategories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(dashboard.ExpenseCategories))
	}
	if dashboard.ExpenseCategories[0].Key != CategoryGroceries {
		t.Fatalf("expected first category supermercado, got %s", dashboard.ExpenseCategories[0].Key)
	}
	if len(dashboard.TopExpenses) == 0 || dashboard.TopExpenses[0].Title != "Supermercado" {
		t.Fatalf("expected top expense Supermercado, got %+v", dashboard.TopExpenses)
	}
	if len(dashboard.RecentTransactions) == 0 || dashboard.RecentTransactions[0].AccountName != "Cuenta sueldo" {
		t.Fatalf("expected recent transactions to include account name, got %+v", dashboard.RecentTransactions)
	}
	if len(dashboard.RecentTransactions) == 0 || dashboard.RecentTransactions[0].Title != "Café" {
		t.Fatalf("expected recent transactions to be independent from active period, got %+v", dashboard.RecentTransactions)
	}
	if len(dashboard.RecentTransactions) != 5 {
		t.Fatalf("expected 5 recent transactions, got %d", len(dashboard.RecentTransactions))
	}
	if dashboard.ActivePeriod.Key != DashboardPeriodMonth {
		t.Fatalf("expected active period month, got %s", dashboard.ActivePeriod.Key)
	}
	if dashboard.Comparison == nil {
		t.Fatal("expected comparison to be returned")
	}
	if dashboard.Comparison.Title != "Mes actual vs anterior" {
		t.Fatalf("expected month comparison title, got %s", dashboard.Comparison.Title)
	}
	if dashboard.Comparison.Previous.IncomeTotal != 2500 {
		t.Fatalf("expected previous month income total 2500, got %v", dashboard.Comparison.Previous.IncomeTotal)
	}
}

func TestGetDashboardBuildsEquivalentComparisonForCustomRange(t *testing.T) {
	repo := &fakeTransactionsRepo{
		transactions: []Transaction{
			{
				ID:        1,
				UserID:    1,
				Title:     "Proyecto",
				Amount:    1500,
				Type:      Income,
				AccountID: 10,
				Category:  "salario",
				CreatedAt: time.Date(2026, 4, 5, 11, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				UserID:    1,
				Title:     "Supermercado",
				Amount:    400,
				Type:      Expense,
				AccountID: 10,
				Category:  "supermercado",
				CreatedAt: time.Date(2026, 4, 6, 18, 0, 0, 0, time.UTC),
			},
			{
				ID:        3,
				UserID:    1,
				Title:     "Proyecto anterior",
				Amount:    1000,
				Type:      Income,
				AccountID: 10,
				Category:  "salario",
				CreatedAt: time.Date(2026, 4, 3, 11, 0, 0, 0, time.UTC),
			},
			{
				ID:        4,
				UserID:    1,
				Title:     "Taxi anterior",
				Amount:    120,
				Type:      Expense,
				AccountID: 10,
				Category:  "taxi",
				CreatedAt: time.Date(2026, 4, 4, 18, 0, 0, 0, time.UTC),
			},
		},
	}
	accountSvc := &fakeAccountsService{
		accounts: []accounts.Account{
			{ID: 10, Name: "Cuenta sueldo"},
		},
	}

	service := NewService(repo, accountSvc, accountSvc)

	dateFrom := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)
	period := DashboardPeriodCustom

	dashboard, err := service.GetDashboard(1, TransactionFilters{
		Period:   &period,
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if dashboard.Comparison == nil {
		t.Fatal("expected custom comparison to be returned")
	}
	if dashboard.Comparison.Title != "Período actual vs anterior" {
		t.Fatalf("expected custom comparison title, got %s", dashboard.Comparison.Title)
	}
	if dashboard.Comparison.Previous.IncomeTotal != 1000 {
		t.Fatalf("expected previous custom income total 1000, got %v", dashboard.Comparison.Previous.IncomeTotal)
	}
	if dashboard.Comparison.Previous.ExpenseTotal != 120 {
		t.Fatalf("expected previous custom expense total 120, got %v", dashboard.Comparison.Previous.ExpenseTotal)
	}
}

func TestGetDashboardCategoryDetailGroupsNormalizedExpenses(t *testing.T) {
	repo := &fakeTransactionsRepo{
		transactions: []Transaction{
			{
				ID:        1,
				UserID:    1,
				Title:     "Restaurante",
				Amount:    500,
				Type:      Expense,
				AccountID: 10,
				Category:  "restaurante",
				CreatedAt: time.Date(2026, 4, 7, 15, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				UserID:    1,
				Title:     "Comida rápida",
				Amount:    200,
				Type:      Expense,
				AccountID: 10,
				Category:  "comida",
				CreatedAt: time.Date(2026, 4, 6, 15, 0, 0, 0, time.UTC),
			},
			{
				ID:        3,
				UserID:    1,
				Title:     "Taxi",
				Amount:    100,
				Type:      Expense,
				AccountID: 10,
				Category:  "taxi",
				CreatedAt: time.Date(2026, 4, 5, 15, 0, 0, 0, time.UTC),
			},
		},
	}
	accountSvc := &fakeAccountsService{
		accounts: []accounts.Account{
			{ID: 10, Name: "Cuenta sueldo"},
		},
	}

	service := NewService(repo, accountSvc, accountSvc)

	dateFrom := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)

	detail, err := service.GetDashboardCategoryDetail(1, TransactionFilters{
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	}, CategoryFood)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if detail.Category.Amount != 700 {
		t.Fatalf("expected category amount 700, got %v", detail.Category.Amount)
	}
	if detail.Category.TransactionsCount != 2 {
		t.Fatalf("expected category transactions count 2, got %d", detail.Category.TransactionsCount)
	}
	if len(detail.Transactions) != 2 {
		t.Fatalf("expected 2 category transactions, got %d", len(detail.Transactions))
	}
	if detail.Transactions[0].CategoryLabel != "Comida" {
		t.Fatalf("expected normalized category label Comida, got %s", detail.Transactions[0].CategoryLabel)
	}
}

func TestCreateRejectsInvalidCategory(t *testing.T) {
	repo := &fakeTransactionsRepo{}
	accountSvc := &fakeAccountsService{
		account: accounts.Account{ID: 5, Name: "Caja", IsActive: true},
		found:   true,
	}

	service := NewService(repo, accountSvc, accountSvc)

	_, err := service.Create(1, CreateTransactionInput{
		Title:     "Compra",
		Amount:    100,
		Type:      Expense,
		AccountID: 5,
		Category:  "mascotas",
	})
	if err != ErrTransactionCategoryInvalid {
		t.Fatalf("expected ErrTransactionCategoryInvalid, got %v", err)
	}
	if repo.createCalled {
		t.Fatal("expected repo.Create not to be called")
	}
}

func TestGetAllNormalizesLegacyCategoriesForFilters(t *testing.T) {
	repo := &fakeTransactionsRepo{
		transactions: []Transaction{
			{
				ID:        1,
				UserID:    1,
				Title:     "Sueldo",
				Amount:    1000,
				Type:      Income,
				AccountID: 10,
				Category:  "salario",
				CreatedAt: time.Date(2026, 4, 7, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				UserID:    1,
				Title:     "Taxi",
				Amount:    100,
				Type:      Expense,
				AccountID: 10,
				Category:  "taxi",
				CreatedAt: time.Date(2026, 4, 7, 11, 0, 0, 0, time.UTC),
			},
		},
	}
	accountSvc := &fakeAccountsService{}
	service := NewService(repo, accountSvc, accountSvc)

	categoryKey := CategorySalary
	transactions, err := service.GetAll(1, TransactionFilters{Category: &categoryKey})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].Category != CategorySalary {
		t.Fatalf("expected normalized category sueldo, got %s", transactions[0].Category)
	}
	if transactions[0].CategoryLabel != "Sueldo" {
		t.Fatalf("expected normalized category label Sueldo, got %s", transactions[0].CategoryLabel)
	}
}

func TestListCategoriesKeepsOtrosAsLastFallback(t *testing.T) {
	service := NewService(&fakeTransactionsRepo{}, &fakeAccountsService{}, &fakeAccountsService{})

	categories := service.ListCategories()
	if len(categories) == 0 {
		t.Fatal("expected categories to be returned")
	}

	lastCategory := categories[len(categories)-1]
	if lastCategory.Key != CategoryOther {
		t.Fatalf("expected last category to be otros, got %s", lastCategory.Key)
	}
}

func pointerToString(value string) *string {
	return &value
}
