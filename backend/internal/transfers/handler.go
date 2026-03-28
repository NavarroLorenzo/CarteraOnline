package transfers

import (
	"cartera-app/backend/internal/auth"
	"cartera-app/backend/internal/shared/httpjson"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	var input CreateTransferInput

	if err := c.ShouldBindJSON(&input); err != nil {
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	err := h.service.CreateTransfer(userID, input)
	if err != nil {
		switch err {
		case ErrSameAccount:
			httpjson.Error(c, 400, "same_account_transfer", err.Error())
		case ErrSourceAccountNotFound:
			httpjson.Error(c, 404, "source_account_not_found", err.Error())
		case ErrDestinationAccountNotFound:
			httpjson.Error(c, 404, "destination_account_not_found", err.Error())
		case ErrSourceAccountInactive:
			httpjson.Error(c, 409, "source_account_inactive", err.Error())
		case ErrDestinationAccountInactive:
			httpjson.Error(c, 409, "destination_account_inactive", err.Error())
		case ErrInsufficientFunds:
			httpjson.Error(c, 409, "insufficient_funds", err.Error())
		case ErrTransferAmountInvalid:
			httpjson.Error(c, 400, "transfer_amount_invalid", err.Error())
		default:
			httpjson.Error(c, 500, "transfer_failed", "No se pudo realizar la transferencia")
		}
		return
	}

	c.JSON(201, gin.H{
		"message": "Transferencia realizada correctamente",
	})
}
