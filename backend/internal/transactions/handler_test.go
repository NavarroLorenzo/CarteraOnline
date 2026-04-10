package transactions

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBuildTransactionFiltersRejectsInvalidDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/transactions?date_from=2026-03-10&date_to=2026-03-01", nil)
	context.Request = request

	_, err := buildTransactionFilters(context)
	if err != ErrInvalidDateRange {
		t.Fatalf("expected ErrInvalidDateRange, got %v", err)
	}
}

func TestBuildTransactionFiltersRejectsInvalidAccountID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/transactions?account_id=abc", nil)
	context.Request = request

	_, err := buildTransactionFilters(context)
	if err != ErrInvalidAccountID {
		t.Fatalf("expected ErrInvalidAccountID, got %v", err)
	}
}

func TestBuildTransactionFiltersRejectsInvalidCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/transactions?category=mascotas", nil)
	context.Request = request

	_, err := buildTransactionFilters(context)
	if err != ErrTransactionCategoryInvalid {
		t.Fatalf("expected ErrTransactionCategoryInvalid, got %v", err)
	}
}

func TestBuildTransactionFiltersParsesDashboardPeriodAndCamelCaseDates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/transactions/dashboard?period=custom&dateFrom=2026-04-01&dateTo=2026-04-07", nil)
	context.Request = request

	filters, err := buildTransactionFilters(context)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filters.Period == nil || *filters.Period != DashboardPeriodCustom {
		t.Fatalf("expected custom dashboard period, got %+v", filters.Period)
	}
	if filters.DateFrom == nil || filters.DateFrom.Format("2006-01-02") != "2026-04-01" {
		t.Fatalf("expected parsed dateFrom, got %+v", filters.DateFrom)
	}
	if filters.DateTo == nil || filters.DateTo.Format("2006-01-02") != "2026-04-07" {
		t.Fatalf("expected parsed dateTo, got %+v", filters.DateTo)
	}
}

func TestBuildTransactionFiltersRejectsCustomPeriodWithoutDates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("GET", "/transactions/dashboard?period=custom&dateFrom=2026-04-01", nil)
	context.Request = request

	_, err := buildTransactionFilters(context)
	if err != ErrMissingDateRange {
		t.Fatalf("expected ErrMissingDateRange, got %v", err)
	}
}
