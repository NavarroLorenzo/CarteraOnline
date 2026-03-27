package transactions

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	transaction, err := h.service.Create(input)
	if err != nil {
		if err == ErrAccountNotFound {
			c.JSON(400, gin.H{
				"error": "La cuenta indicada no existe",
			})
			return
		}

		c.JSON(500, gin.H{
			"error": "No se pudo crear la transacción",
		})
		return
	}

	c.JSON(201, transaction)
}

func (h *Handler) GetAll(c *gin.Context) {
	transactions := h.service.GetAll()
	c.JSON(200, transactions)
}

func (h *Handler) GetBalance(c *gin.Context) {
	balance := h.service.GetBalance()

	c.JSON(200, gin.H{
		"balance": balance,
	})
}

func (h *Handler) GetBalanceByAccount(c *gin.Context) {
	accounts, total := h.service.GetBalanceByAccountDetailed()

	c.JSON(200, gin.H{
		"accounts": accounts,
		"total":    total,
	})
}
