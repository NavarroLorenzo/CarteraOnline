package accounts

import (
	"cartera-app/backend/internal/auth"
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
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	account, err := h.service.Create(userID, CreateAccountInput{
		Name: input.Name,
		Type: input.Type,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear la cuenta"})
		return
	}

	if input.InitialAmount > 0 {
		err = h.initialBalanceCreator.CreateInitialBalance(userID, account.ID, input.InitialAmount)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "La cuenta se creó, pero falló la carga del saldo inicial",
			})
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
		c.JSON(500, gin.H{"error": "No se pudieron obtener las cuentas"})
		return
	}

	c.JSON(200, accounts)
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "id inválido"})
		return
	}

	var input UpdateAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	account, err := h.service.Update(userID, id, input)
	if err != nil {
		switch err {
		case ErrAccountNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo actualizar la cuenta"})
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
		c.JSON(400, gin.H{"error": "id inválido"})
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		switch err {
		case ErrAccountNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo eliminar la cuenta"})
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Cuenta eliminada correctamente",
		"warning": "También se eliminaron las transacciones asociadas a la cuenta",
	})
}
