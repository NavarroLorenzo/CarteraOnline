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
