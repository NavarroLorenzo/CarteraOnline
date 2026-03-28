package transfers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"cartera-app/backend/internal/accounts"
	"cartera-app/backend/internal/shared/normalize"
)

var (
	ErrSameAccount                = errors.New("no se puede transferir a la misma cuenta")
	ErrSourceAccountNotFound      = errors.New("la cuenta origen no existe")
	ErrDestinationAccountNotFound = errors.New("la cuenta destino no existe")
	ErrSourceAccountInactive      = errors.New("la cuenta origen está inactiva")
	ErrDestinationAccountInactive = errors.New("la cuenta destino está inactiva")
	ErrInsufficientFunds          = errors.New("saldo insuficiente en la cuenta origen")
	ErrTransferAmountInvalid      = errors.New("el monto debe ser mayor a cero")
)

type AccountProvider interface {
	GetByID(userID, id int64) (accounts.Account, bool, error)
}

type Service struct {
	accountProvider AccountProvider
	repo            Repository
}

func NewService(accountProvider AccountProvider, repo Repository) *Service {
	return &Service{
		accountProvider: accountProvider,
		repo:            repo,
	}
}

func (s *Service) CreateTransfer(userID int64, input CreateTransferInput) error {
	amount, err := normalize.Money(input.Amount)
	if err != nil {
		return ErrTransferAmountInvalid
	}

	if input.FromAccountID == input.ToAccountID {
		return ErrSameAccount
	}

	fromAccount, found, err := s.accountProvider.GetByID(userID, input.FromAccountID)
	if err != nil {
		return err
	}
	if !found {
		return ErrSourceAccountNotFound
	}
	if !fromAccount.IsActive {
		return ErrSourceAccountInactive
	}

	toAccount, found, err := s.accountProvider.GetByID(userID, input.ToAccountID)
	if err != nil {
		return err
	}
	if !found {
		return ErrDestinationAccountNotFound
	}
	if !toAccount.IsActive {
		return ErrDestinationAccountInactive
	}

	expenseDescription := "Transferencia enviada a cuenta: " + toAccount.Name
	incomeDescription := "Transferencia recibida desde cuenta: " + fromAccount.Name

	if userDescription := normalize.Optional(input.Description); userDescription != "" {
		expenseDescription += ". Detalle: " + userDescription
		incomeDescription += ". Detalle: " + userDescription
	}

	transferID, err := generateTransferID()
	if err != nil {
		return err
	}

	return s.repo.ExecuteTransfer(userID, CreateTransferOperation{
		FromAccountID:      input.FromAccountID,
		ToAccountID:        input.ToAccountID,
		Amount:             amount,
		TransferID:         transferID,
		ExpenseTitle:       "Transferencia enviada",
		ExpenseDescription: expenseDescription,
		IncomeTitle:        "Transferencia recibida",
		IncomeDescription:  incomeDescription,
	})
}

func generateTransferID() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return "tr_" + hex.EncodeToString(buffer), nil
}
