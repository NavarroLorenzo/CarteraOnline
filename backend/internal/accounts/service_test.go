package accounts

import "testing"

type fakeAccountsRepository struct {
	account      Account
	found        bool
	deleteCalled bool
}

func (f *fakeAccountsRepository) Create(input CreateAccountInput) (Account, error) {
	return Account{}, nil
}

func (f *fakeAccountsRepository) GetAll(userID int64) ([]Account, error) {
	return nil, nil
}

func (f *fakeAccountsRepository) GetByID(userID, id int64) (Account, bool, error) {
	return f.account, f.found, nil
}

func (f *fakeAccountsRepository) ExistsByID(userID, id int64) (bool, error) {
	return f.found, nil
}

func (f *fakeAccountsRepository) ExistsActiveByNormalizedName(userID int64, nameNormalized string, excludeID *int64) (bool, error) {
	return false, nil
}

func (f *fakeAccountsRepository) Update(userID, id int64, input UpdateAccountInput) (Account, bool, error) {
	return Account{}, false, nil
}

func (f *fakeAccountsRepository) Delete(userID, id int64) (bool, error) {
	f.deleteCalled = true
	return true, nil
}

func TestDeleteRemovesAccountEvenIfItHasTransactions(t *testing.T) {
	repo := &fakeAccountsRepository{
		account: Account{ID: 10, Name: "Caja"},
		found:   true,
	}

	service := NewService(repo)

	err := service.Delete(1, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.deleteCalled {
		t.Fatal("expected repository delete to be called")
	}
}
