package transactions

import (
	"cartera-app/backend/internal/auth"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidAccountID = errors.New("account_id debe ser numérico")
	ErrInvalidType      = errors.New("type debe ser income o expense")
	ErrInvalidDateFrom  = errors.New("date_from debe tener formato YYYY-MM-DD")
	ErrInvalidDateTo    = errors.New("date_to debe tener formato YYYY-MM-DD")
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	var input CreateTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	transaction, err := h.service.Create(userID, input)
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
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	filters, err := buildTransactionFilters(c)
	if err != nil {
		h.handleFilterError(c, err)
		return
	}

	transactions, err := h.service.GetAll(userID, filters)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "No se pudieron obtener las transacciones",
		})
		return
	}

	c.JSON(200, transactions)
}

func (h *Handler) GetBalance(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo calcular el balance"})
		return
	}

	c.JSON(200, gin.H{"balance": balance})
}

func (h *Handler) GetBalanceByAccount(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	accounts, total, err := h.service.GetBalanceByAccountDetailed(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo obtener el balance por cuenta"})
		return
	}

	c.JSON(200, gin.H{
		"accounts": accounts,
		"total":    total,
	})
}

func buildTransactionFilters(c *gin.Context) (TransactionFilters, error) {
	var filters TransactionFilters

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			return filters, ErrInvalidAccountID
		}
		filters.AccountID = &accountID
	}

	if typeStr := c.Query("type"); typeStr != "" {
		transactionType := TransactionType(typeStr)

		if transactionType != Income && transactionType != Expense {
			return filters, ErrInvalidType
		}

		filters.Type = &transactionType
	}

	if categoryStr := c.Query("category"); categoryStr != "" {
		filters.Category = &categoryStr
	}

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err != nil {
			return filters, ErrInvalidDateFrom
		}
		filters.DateFrom = &dateFrom
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			return filters, ErrInvalidDateTo
		}

		dateTo = dateTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		filters.DateTo = &dateTo
	}

	return filters, nil
}

func (h *Handler) handleFilterError(c *gin.Context, err error) {
	switch err {
	case ErrInvalidAccountID:
		c.JSON(400, gin.H{"error": ErrInvalidAccountID.Error()})
	case ErrInvalidType:
		c.JSON(400, gin.H{"error": ErrInvalidType.Error()})
	case ErrInvalidDateFrom:
		c.JSON(400, gin.H{"error": ErrInvalidDateFrom.Error()})
	case ErrInvalidDateTo:
		c.JSON(400, gin.H{"error": ErrInvalidDateTo.Error()})
	default:
		c.JSON(400, gin.H{"error": "Filtros inválidos"})
	}
}

func (h *Handler) GetSummary(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	filters, err := buildTransactionFilters(c)
	if err != nil {
		h.handleFilterError(c, err)
		return
	}

	summary, err := h.service.GetSummary(userID, filters)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "No se pudo obtener el resumen de transacciones",
		})
		return
	}

	c.JSON(200, summary)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "id inválido"})
		return
	}

	transaction, err := h.service.GetByID(userID, id)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo obtener la transacción"})
		}
		return
	}

	c.JSON(200, transaction)
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

	var input UpdateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	transaction, err := h.service.Update(userID, id, input)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		case ErrAccountNotFound:
			c.JSON(400, gin.H{"error": "La cuenta indicada no existe"})
		case ErrCannotUpdateTransfer:
			c.JSON(400, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo actualizar la transacción"})
		}
		return
	}

	c.JSON(200, transaction)
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

	err = h.service.Delete(userID, id)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo eliminar la transacción"})
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Transacción eliminada correctamente",
	})
}
