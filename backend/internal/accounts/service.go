package accounts

import "errors"

var ErrAccountNotFound = errors.New("la cuenta no existe")

type Service interface {
	Create(userID int64, input CreateAccountInput) (Account, error)
	GetAll(userID int64) ([]Account, error)
	GetByID(userID, id int64) (Account, bool, error)
	ExistsByID(userID, id int64) (bool, error)
	Update(userID, id int64, input UpdateAccountInput) (Account, error)
	Delete(userID, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(userID int64, input CreateAccountInput) (Account, error) {
	input.UserID = userID
	return s.repo.Create(input)
}

func (s *service) GetAll(userID int64) ([]Account, error) {
	return s.repo.GetAll(userID)
}

func (s *service) GetByID(userID, id int64) (Account, bool, error) {
	return s.repo.GetByID(userID, id)
}

func (s *service) ExistsByID(userID, id int64) (bool, error) {
	return s.repo.ExistsByID(userID, id)
}

func (s *service) Update(userID, id int64, input UpdateAccountInput) (Account, error) {
	account, found, err := s.repo.Update(userID, id, input)
	if err != nil {
		return Account{}, err
	}
	if !found {
		return Account{}, ErrAccountNotFound
	}

	return account, nil
}

func (s *service) Delete(userID, id int64) error {
	deleted, err := s.repo.Delete(userID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrAccountNotFound
	}

	return nil
}
