package accounts

import "github.com/gin-gonic/gin"

type InitialBalanceCreator interface {
	CreateInitialBalance(accountID int64, amount float64) error
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
	var input CreateAccountInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	account, err := h.service.Create(CreateAccountInput{
		Name: input.Name,
		Type: input.Type,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear la cuenta"})
		return
	}

	if input.InitialAmount > 0 {
		err = h.initialBalanceCreator.CreateInitialBalance(account.ID, input.InitialAmount)
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
	accounts, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudieron obtener las cuentas"})
		return
	}

	c.JSON(200, accounts)
}
