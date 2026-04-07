package transactions

import (
	"cartera-app/backend/internal/auth"
	"cartera-app/backend/internal/shared/httpjson"
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
	ErrInvalidDateRange = errors.New("date_from no puede ser mayor que date_to")
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
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	transaction, err := h.service.Create(userID, input)
	if err != nil {
		switch err {
		case ErrAccountNotFound:
			httpjson.Error(c, 404, "account_not_found", err.Error())
		case ErrAccountInactive:
			httpjson.Error(c, 409, "account_inactive", err.Error())
		case ErrTransactionTitleRequired:
			httpjson.Error(c, 400, "transaction_title_required", err.Error())
		case ErrTransactionCategoryRequired:
			httpjson.Error(c, 400, "transaction_category_required", err.Error())
		case ErrTransactionAmountInvalid:
			httpjson.Error(c, 400, "transaction_amount_invalid", err.Error())
		case ErrTransactionTypeInvalid:
			httpjson.Error(c, 400, "transaction_type_invalid", err.Error())
		default:
			httpjson.Error(c, 500, "transaction_create_failed", "No se pudo crear la transacción")
		}
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
		httpjson.Error(c, 500, "transactions_list_failed", "No se pudieron obtener las transacciones")
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
		httpjson.Error(c, 500, "balance_failed", "No se pudo calcular el balance")
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
		httpjson.Error(c, 500, "balance_by_account_failed", "No se pudo obtener el balance por cuenta")
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

	if filters.DateFrom != nil && filters.DateTo != nil && filters.DateFrom.After(*filters.DateTo) {
		return filters, ErrInvalidDateRange
	}

	return filters, nil
}

func (h *Handler) handleFilterError(c *gin.Context, err error) {
	switch err {
	case ErrInvalidAccountID:
		httpjson.Error(c, 400, "invalid_account_id", ErrInvalidAccountID.Error())
	case ErrInvalidType:
		httpjson.Error(c, 400, "invalid_transaction_type", ErrInvalidType.Error())
	case ErrInvalidDateFrom:
		httpjson.Error(c, 400, "invalid_date_from", ErrInvalidDateFrom.Error())
	case ErrInvalidDateTo:
		httpjson.Error(c, 400, "invalid_date_to", ErrInvalidDateTo.Error())
	case ErrInvalidDateRange:
		httpjson.Error(c, 400, "invalid_date_range", ErrInvalidDateRange.Error())
	default:
		httpjson.Error(c, 400, "invalid_filters", "Filtros inválidos")
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
		httpjson.Error(c, 500, "summary_failed", "No se pudo obtener el resumen de transacciones")
		return
	}

	c.JSON(200, summary)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	filters, err := buildTransactionFilters(c)
	if err != nil {
		h.handleFilterError(c, err)
		return
	}

	dashboard, err := h.service.GetDashboard(userID, filters)
	if err != nil {
		httpjson.Error(c, 500, "dashboard_failed", "No se pudo obtener la analítica del dashboard")
		return
	}

	c.JSON(200, dashboard)
}

func (h *Handler) GetDashboardCategoryDetail(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	filters, err := buildTransactionFilters(c)
	if err != nil {
		h.handleFilterError(c, err)
		return
	}

	categoryKey := c.Query("category_key")
	detail, err := h.service.GetDashboardCategoryDetail(userID, filters, categoryKey)
	if err != nil {
		switch err {
		case ErrDashboardCategoryRequired:
			httpjson.Error(c, 400, "dashboard_category_required", err.Error())
		case ErrInvalidDashboardCategory:
			httpjson.Error(c, 400, "dashboard_category_invalid", err.Error())
		default:
			httpjson.Error(c, 500, "dashboard_category_failed", "No se pudo obtener el detalle de la categoría")
		}
		return
	}

	c.JSON(200, detail)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpjson.Error(c, 400, "invalid_transaction_id", "id inválido")
		return
	}

	transaction, err := h.service.GetByID(userID, id)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			httpjson.Error(c, 404, "transaction_not_found", err.Error())
		default:
			httpjson.Error(c, 500, "transaction_get_failed", "No se pudo obtener la transacción")
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
		httpjson.Error(c, 400, "invalid_transaction_id", "id inválido")
		return
	}

	var input UpdateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	transaction, err := h.service.Update(userID, id, input)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			httpjson.Error(c, 404, "transaction_not_found", err.Error())
		case ErrAccountNotFound:
			httpjson.Error(c, 404, "account_not_found", err.Error())
		case ErrAccountInactive:
			httpjson.Error(c, 409, "account_inactive", err.Error())
		case ErrCannotUpdateTransfer:
			httpjson.Error(c, 409, "transfer_update_blocked", err.Error())
		case ErrTransactionTitleRequired:
			httpjson.Error(c, 400, "transaction_title_required", err.Error())
		case ErrTransactionCategoryRequired:
			httpjson.Error(c, 400, "transaction_category_required", err.Error())
		case ErrTransactionAmountInvalid:
			httpjson.Error(c, 400, "transaction_amount_invalid", err.Error())
		case ErrTransactionTypeInvalid:
			httpjson.Error(c, 400, "transaction_type_invalid", err.Error())
		default:
			httpjson.Error(c, 500, "transaction_update_failed", "No se pudo actualizar la transacción")
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
		httpjson.Error(c, 400, "invalid_transaction_id", "id inválido")
		return
	}

	err = h.service.Delete(userID, id)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			httpjson.Error(c, 404, "transaction_not_found", err.Error())
		default:
			httpjson.Error(c, 500, "transaction_delete_failed", "No se pudo eliminar la transacción")
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Transacción eliminada correctamente",
	})
}
