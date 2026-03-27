package transactions

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
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
			c.JSON(400, gin.H{"error": "La cuenta indicada no existe"})
			return
		}
		c.JSON(500, gin.H{"error": "No se pudo crear la transacción"})
		return
	}

	c.JSON(201, transaction)
}

func (h *Handler) GetAll(c *gin.Context) {
	var filters TransactionFilters

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "account_id debe ser numérico",
			})
			return
		}
		filters.AccountID = &accountID
	}

	if typeStr := c.Query("type"); typeStr != "" {
		transactionType := TransactionType(typeStr)

		if transactionType != Income && transactionType != Expense {
			c.JSON(400, gin.H{
				"error": "type debe ser income o expense",
			})
			return
		}

		filters.Type = &transactionType
	}

	if categoryStr := c.Query("category"); categoryStr != "" {
		filters.Category = &categoryStr
	}

	transactions, err := h.service.GetAll(filters)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "No se pudieron obtener las transacciones",
		})
		return
	}

	c.JSON(200, transactions)
}

func (h *Handler) GetBalance(c *gin.Context) {
	balance, err := h.service.GetBalance()
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo calcular el balance"})
		return
	}

	c.JSON(200, gin.H{"balance": balance})
}

func (h *Handler) GetBalanceByAccount(c *gin.Context) {
	accounts, total, err := h.service.GetBalanceByAccountDetailed()
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo obtener el balance por cuenta"})
		return
	}

	c.JSON(200, gin.H{
		"accounts": accounts,
		"total":    total,
	})
}
