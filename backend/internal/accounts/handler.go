package accounts

import (
	"cartera-app/backend/internal/auth"
	"cartera-app/backend/internal/shared/httpjson"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InitialBalanceCreator interface {
	CreateInitialBalance(userID, accountID int64, amount float64) error
}

type Handler struct {
	service               Service
	initialBalanceCreator InitialBalanceCreator
}

func NewHandler(service Service, initialBalanceCreator InitialBalanceCreator) *Handler {
	return &Handler{
		service:               service,
		initialBalanceCreator: initialBalanceCreator,
	}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	var input CreateAccountInput

	if err := c.ShouldBindJSON(&input); err != nil {
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	account, err := h.service.Create(userID, input)
	if err != nil {
		switch err {
		case ErrAccountNameRequired:
			httpjson.Error(c, 400, "account_name_required", err.Error())
		case ErrInvalidAccountType:
			httpjson.Error(c, 400, "invalid_account_type", err.Error())
		case ErrDuplicateActiveAccount:
			httpjson.Error(c, 409, "duplicate_account_name", err.Error())
		default:
			httpjson.Error(c, 500, "account_create_failed", "No se pudo crear la cuenta")
		}
		return
	}

	if input.InitialAmount > 0 {
		err = h.initialBalanceCreator.CreateInitialBalance(userID, account.ID, input.InitialAmount)
		if err != nil {
			_ = h.service.Delete(userID, account.ID)
			httpjson.Error(c, 500, "initial_balance_failed", "No se pudo crear la cuenta con su saldo inicial")
			return
		}
	}

	c.JSON(201, account)
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	accounts, err := h.service.GetAll(userID)
	if err != nil {
		httpjson.Error(c, 500, "accounts_list_failed", "No se pudieron obtener las cuentas")
		return
	}

	c.JSON(200, accounts)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpjson.Error(c, 400, "invalid_account_id", "id inválido")
		return
	}

	account, found, err := h.service.GetByID(userID, id)
	if err != nil {
		httpjson.Error(c, 500, "account_get_failed", "No se pudo obtener la cuenta")
		return
	}
	if !found {
		httpjson.Error(c, 404, "account_not_found", ErrAccountNotFound.Error())
		return
	}

	c.JSON(200, account)
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpjson.Error(c, 400, "invalid_account_id", "id inválido")
		return
	}

	var input UpdateAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	account, err := h.service.Update(userID, id, input)
	if err != nil {
		switch err {
		case ErrAccountNotFound:
			httpjson.Error(c, 404, "account_not_found", err.Error())
		case ErrAccountNameRequired:
			httpjson.Error(c, 400, "account_name_required", err.Error())
		case ErrInvalidAccountType:
			httpjson.Error(c, 400, "invalid_account_type", err.Error())
		case ErrDuplicateActiveAccount:
			httpjson.Error(c, 409, "duplicate_account_name", err.Error())
		default:
			httpjson.Error(c, 500, "account_update_failed", "No se pudo actualizar la cuenta")
		}
		return
	}

	c.JSON(200, account)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpjson.Error(c, 400, "invalid_account_id", "id inválido")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		switch err {
		case ErrAccountNotFound:
			httpjson.Error(c, 404, "account_not_found", err.Error())
		default:
			httpjson.Error(c, 500, "account_delete_failed", "No se pudo eliminar la cuenta")
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Cuenta eliminada correctamente junto con las transacciones de esa cuenta",
	})
}
