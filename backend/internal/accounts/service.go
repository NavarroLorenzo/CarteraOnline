package accounts

import (
	"errors"
	"strings"

	"cartera-app/backend/internal/shared/normalize"
)

var ErrAccountNotFound = errors.New("la cuenta no existe")
var ErrAccountNameRequired = errors.New("el nombre de la cuenta es obligatorio")
var ErrInvalidAccountType = errors.New("el tipo de cuenta no es válido")
var ErrDuplicateActiveAccount = errors.New("ya existe una cuenta activa con ese nombre")

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
	name, nameNormalized, err := normalizeAccountName(input.Name)
	if err != nil {
		return Account{}, err
	}

	accountType, err := normalizeAccountType(input.Type)
	if err != nil {
		return Account{}, err
	}

	exists, err := s.repo.ExistsActiveByNormalizedName(userID, nameNormalized, nil)
	if err != nil {
		return Account{}, err
	}
	if exists {
		return Account{}, ErrDuplicateActiveAccount
	}

	input.UserID = userID
	input.Name = name
	input.NameNormalized = nameNormalized
	input.Type = accountType
	input.IsActive = true
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
	current, found, err := s.repo.GetByID(userID, id)
	if err != nil {
		return Account{}, err
	}
	if !found {
		return Account{}, ErrAccountNotFound
	}

	name, nameNormalized, err := normalizeAccountName(input.Name)
	if err != nil {
		return Account{}, err
	}

	accountType, err := normalizeAccountType(input.Type)
	if err != nil {
		return Account{}, err
	}

	isActive := current.IsActive
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	if isActive {
		exists, err := s.repo.ExistsActiveByNormalizedName(userID, nameNormalized, &id)
		if err != nil {
			return Account{}, err
		}
		if exists {
			return Account{}, ErrDuplicateActiveAccount
		}
	}

	account, found, err := s.repo.Update(userID, id, UpdateAccountInput{
		Name:           name,
		NameNormalized: nameNormalized,
		Type:           accountType,
		IsActive:       &isActive,
	})
	if err != nil {
		return Account{}, err
	}
	if !found {
		return Account{}, ErrAccountNotFound
	}

	return account, nil
}

func (s *service) Delete(userID, id int64) error {
	_, found, err := s.repo.GetByID(userID, id)
	if err != nil {
		return err
	}
	if !found {
		return ErrAccountNotFound
	}

	deleted, err := s.repo.Delete(userID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrAccountNotFound
	}

	return nil
}

func normalizeAccountName(value string) (string, string, error) {
	name := normalize.Optional(value)
	if name == "" {
		return "", "", ErrAccountNameRequired
	}

	return name, normalize.LowerKey(name), nil
}

func normalizeAccountType(value string) (string, error) {
	accountType := strings.ToLower(strings.TrimSpace(value))
	if _, ok := AllowedTypes[accountType]; !ok {
		return "", ErrInvalidAccountType
	}

	return accountType, nil
}
