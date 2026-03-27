package transfers

import (
	"errors"
	"fmt"
	"time"

	"cartera-app/backend/internal/accounts"
	"cartera-app/backend/internal/transactions"
)

var (
	ErrSameAccount     = errors.New("no se puede transferir a la misma cuenta")
	ErrAccountNotFound = errors.New("una de las cuentas no existe")
)

type AccountProvider interface {
	GetByID(id int64) (accounts.Account, bool, error)
}

type TransactionCreator interface {
	Create(input transactions.CreateTransactionInput) (transactions.Transaction, error)
}

type Service struct {
	accountProvider AccountProvider
	transactionSvc  TransactionCreator
}

func NewService(accountProvider AccountProvider, transactionSvc TransactionCreator) *Service {
	return &Service{
		accountProvider: accountProvider,
		transactionSvc:  transactionSvc,
	}
}

func (s *Service) CreateTransfer(input CreateTransferInput) error {
	if input.FromAccountID == input.ToAccountID {
		return ErrSameAccount
	}

	fromAccount, found, err := s.accountProvider.GetByID(input.FromAccountID)
	if err != nil || !found {
		return ErrAccountNotFound
	}

	toAccount, found, err := s.accountProvider.GetByID(input.ToAccountID)
	if err != nil || !found {
		return ErrAccountNotFound
	}

	transferID := fmt.Sprintf("transfer_%d", time.Now().UnixNano())

	expenseDescription := "Transferido a cuenta: " + toAccount.Name
	incomeDescription := "Transferido desde cuenta: " + fromAccount.Name

	if input.Description != "" {
		expenseDescription += ". " + input.Description
		incomeDescription += ". " + input.Description
	}

	_, err = s.transactionSvc.Create(transactions.CreateTransactionInput{
		Title:       "Transferencia enviada",
		Amount:      input.Amount,
		Type:        transactions.Expense,
		AccountID:   input.FromAccountID,
		Category:    "transfer",
		Description: expenseDescription,
		TransferID:  &transferID,
	})
	if err != nil {
		return err
	}

	_, err = s.transactionSvc.Create(transactions.CreateTransactionInput{
		Title:       "Transferencia recibida",
		Amount:      input.Amount,
		Type:        transactions.Income,
		AccountID:   input.ToAccountID,
		Category:    "transfer",
		Description: incomeDescription,
		TransferID:  &transferID,
	})
	if err != nil {
		return err
	}

	return nil
}
