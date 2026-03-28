package transactions

import (
	"testing"

	"cartera-app/backend/internal/accounts"
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
}

func (f *fakeTransactionsRepo) Create(userID int64, input CreateTransactionInput) (Transaction, error) {
	f.createCalled = true
	return Transaction{}, nil
}

func (f *fakeTransactionsRepo) GetAll(userID int64, filters TransactionFilters) ([]Transaction, error) {
	return nil, nil
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
	account accounts.Account
	found   bool
}

func (f *fakeAccountsService) Create(userID int64, input accounts.CreateAccountInput) (accounts.Account, error) {
	return accounts.Account{}, nil
}

func (f *fakeAccountsService) GetAll(userID int64) ([]accounts.Account, error) {
	return nil, nil
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
		Category:  "hogar",
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
