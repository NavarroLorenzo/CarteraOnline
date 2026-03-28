package transfers

import (
	"testing"

	"cartera-app/backend/internal/accounts"
)

type fakeTransferAccountProvider struct {
	accounts map[int64]accounts.Account
}

func (f *fakeTransferAccountProvider) GetByID(userID, id int64) (accounts.Account, bool, error) {
	account, found := f.accounts[id]
	return account, found, nil
}

type fakeTransferRepository struct {
	called bool
	err    error
}

func (f *fakeTransferRepository) ExecuteTransfer(userID int64, input CreateTransferOperation) error {
	f.called = true
	return f.err
}

func TestCreateTransferRejectsInactiveSourceAccount(t *testing.T) {
	repo := &fakeTransferRepository{}
	provider := &fakeTransferAccountProvider{
		accounts: map[int64]accounts.Account{
			1: {ID: 1, Name: "Caja", IsActive: false},
			2: {ID: 2, Name: "Banco", IsActive: true},
		},
	}

	service := NewService(provider, repo)

	err := service.CreateTransfer(1, CreateTransferInput{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        100,
		Description:   "Prueba",
	})
	if err != ErrSourceAccountInactive {
		t.Fatalf("expected ErrSourceAccountInactive, got %v", err)
	}
	if repo.called {
		t.Fatal("expected transfer repository not to be called")
	}
}

func TestCreateTransferPropagatesInsufficientFunds(t *testing.T) {
	repo := &fakeTransferRepository{err: ErrInsufficientFunds}
	provider := &fakeTransferAccountProvider{
		accounts: map[int64]accounts.Account{
			1: {ID: 1, Name: "Caja", IsActive: true},
			2: {ID: 2, Name: "Banco", IsActive: true},
		},
	}

	service := NewService(provider, repo)

	err := service.CreateTransfer(1, CreateTransferInput{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        5000,
		Description:   "Saldo",
	})
	if err != ErrInsufficientFunds {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}
	if !repo.called {
		t.Fatal("expected transfer repository to be called")
	}
}
