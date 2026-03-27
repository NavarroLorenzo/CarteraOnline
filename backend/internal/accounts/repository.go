package accounts

import "time"

type Repository interface {
	Create(input CreateAccountInput) Account
	GetAll() []Account
	GetByID(id int64) (Account, bool)
	ExistsByID(id int64) bool
}

type InMemoryRepository struct {
	data   []Account
	nextID int64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data:   []Account{},
		nextID: 1,
	}
}

func (r *InMemoryRepository) Create(input CreateAccountInput) Account {
	account := Account{
		ID:        r.nextID,
		Name:      input.Name,
		Type:      input.Type,
		CreatedAt: time.Now(),
	}

	r.data = append(r.data, account)
	r.nextID++

	return account
}

func (r *InMemoryRepository) GetAll() []Account {
	return r.data
}

func (r *InMemoryRepository) GetByID(id int64) (Account, bool) {
	for _, account := range r.data {
		if account.ID == id {
			return account, true
		}
	}

	return Account{}, false
}

func (r *InMemoryRepository) ExistsByID(id int64) bool {
	for _, account := range r.data {
		if account.ID == id {
			return true
		}
	}

	return false
}
