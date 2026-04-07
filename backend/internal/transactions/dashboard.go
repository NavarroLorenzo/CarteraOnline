package transactions

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var (
	ErrDashboardCategoryRequired = errors.New("category_key es obligatorio")
	ErrInvalidDashboardCategory  = errors.New("category_key inválido")
)

const (
	DashboardCategoryFood      = "comida"
	DashboardCategoryTransport = "transporte"
	DashboardCategoryLeisure   = "ocio"
	DashboardCategoryServices  = "servicios"
	DashboardCategoryOther     = "otros"
)

type DashboardSummary struct {
	IncomeTotal       float64 `json:"income_total"`
	ExpenseTotal      float64 `json:"expense_total"`
	NetBalance        float64 `json:"net_balance"`
	TransactionsCount int64   `json:"transactions_count"`
}

type DashboardComparison struct {
	CurrentMonth     DashboardSummary `json:"current_month"`
	PreviousMonth    DashboardSummary `json:"previous_month"`
	IncomeChangePct  float64          `json:"income_change_pct"`
	ExpenseChangePct float64          `json:"expense_change_pct"`
}

type DashboardTrendPoint struct {
	Label        string    `json:"label"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IncomeTotal  float64   `json:"income_total"`
	ExpenseTotal float64   `json:"expense_total"`
}

type DashboardCategorySummary struct {
	Key               string  `json:"key"`
	Label             string  `json:"label"`
	Amount            float64 `json:"amount"`
	Percentage        float64 `json:"percentage"`
	TransactionsCount int64   `json:"transactions_count"`
}

type DashboardTransactionItem struct {
	ID            int64           `json:"id"`
	Title         string          `json:"title"`
	Amount        float64         `json:"amount"`
	Type          TransactionType `json:"type"`
	AccountID     int64           `json:"account_id"`
	AccountName   string          `json:"account_name"`
	Category      string          `json:"category"`
	CategoryLabel string          `json:"category_label"`
	Description   string          `json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
}

type DashboardAnalytics struct {
	PeriodSummary      DashboardSummary           `json:"period_summary"`
	DaySummary         DashboardSummary           `json:"day_summary"`
	MonthSummary       DashboardSummary           `json:"month_summary"`
	Comparison         DashboardComparison        `json:"comparison"`
	TrendInterval      string                     `json:"trend_interval"`
	Trend              []DashboardTrendPoint      `json:"trend"`
	ExpenseCategories  []DashboardCategorySummary `json:"expense_categories"`
	TopExpenses        []DashboardTransactionItem `json:"top_expenses"`
	RecentTransactions []DashboardTransactionItem `json:"recent_transactions"`
}

type DashboardCategoryDetail struct {
	Category      DashboardCategorySummary   `json:"category"`
	TrendInterval string                     `json:"trend_interval"`
	Trend         []DashboardTrendPoint      `json:"trend"`
	Transactions  []DashboardTransactionItem `json:"transactions"`
}

type dashboardCategoryDefinition struct {
	Key      string
	Label    string
	Keywords []string
}

var dashboardCategories = []dashboardCategoryDefinition{
	{
		Key:      DashboardCategoryFood,
		Label:    "Comida",
		Keywords: []string{"comida", "super", "supermercado", "grocery", "restaurante", "restaurant", "almacen", "delivery", "cafe"},
	},
	{
		Key:      DashboardCategoryTransport,
		Label:    "Transporte",
		Keywords: []string{"transporte", "nafta", "combustible", "uber", "cabify", "taxi", "sube", "colectivo", "tren", "peaje", "estacionamiento"},
	},
	{
		Key:      DashboardCategoryLeisure,
		Label:    "Ocio",
		Keywords: []string{"ocio", "entretenimiento", "cine", "netflix", "spotify", "juego", "gaming", "salida", "streaming"},
	},
	{
		Key:      DashboardCategoryServices,
		Label:    "Servicios",
		Keywords: []string{"servicio", "servicios", "internet", "telefono", "teléfono", "luz", "agua", "gas", "alquiler", "expensa", "seguro"},
	},
	{
		Key:   DashboardCategoryOther,
		Label: "Otros",
	},
}

func (s *service) GetDashboard(userID int64, filters TransactionFilters) (DashboardAnalytics, error) {
	rangeStart, rangeEnd := resolveDashboardRange(filters)

	periodTransactions, err := s.repo.GetAll(userID, buildDashboardFilters(filters, rangeStart, rangeEnd))
	if err != nil {
		return DashboardAnalytics{}, err
	}

	accountNames, err := s.loadAccountNames(userID)
	if err != nil {
		return DashboardAnalytics{}, err
	}

	analyticsTransactions := filterDashboardAnalyticsTransactions(periodTransactions)
	dayStart := startOfDay(rangeEnd)
	dayEnd := endOfDay(rangeEnd)

	monthStart := startOfMonth(rangeEnd)
	currentMonthTransactions, err := s.repo.GetAll(userID, buildDashboardFilters(filters, monthStart, rangeEnd))
	if err != nil {
		return DashboardAnalytics{}, err
	}

	previousMonthStart, previousMonthEnd := comparablePreviousMonthRange(rangeEnd)
	previousMonthTransactions, err := s.repo.GetAll(userID, buildDashboardFilters(filters, previousMonthStart, previousMonthEnd))
	if err != nil {
		return DashboardAnalytics{}, err
	}

	currentMonthSummary := summarizeTransactions(filterDashboardAnalyticsTransactions(currentMonthTransactions))
	previousMonthSummary := summarizeTransactions(filterDashboardAnalyticsTransactions(previousMonthTransactions))
	trendInterval := selectTrendInterval(rangeStart, rangeEnd)

	return DashboardAnalytics{
		PeriodSummary:      summarizeTransactions(analyticsTransactions),
		DaySummary:         summarizeTransactions(filterTransactionsInRange(analyticsTransactions, dayStart, dayEnd)),
		MonthSummary:       currentMonthSummary,
		Comparison:         buildDashboardComparison(currentMonthSummary, previousMonthSummary),
		TrendInterval:      trendInterval,
		Trend:              buildTrendPoints(analyticsTransactions, rangeStart, rangeEnd, trendInterval),
		ExpenseCategories:  buildCategoryBreakdown(analyticsTransactions),
		TopExpenses:        buildTopExpenses(analyticsTransactions, accountNames),
		RecentTransactions: buildRecentTransactions(periodTransactions, accountNames),
	}, nil
}

func (s *service) GetDashboardCategoryDetail(userID int64, filters TransactionFilters, categoryKey string) (DashboardCategoryDetail, error) {
	normalizedKey, categoryLabel, err := resolveDashboardCategory(categoryKey)
	if err != nil {
		return DashboardCategoryDetail{}, err
	}

	rangeStart, rangeEnd := resolveDashboardRange(filters)
	periodTransactions, err := s.repo.GetAll(userID, buildDashboardFilters(filters, rangeStart, rangeEnd))
	if err != nil {
		return DashboardCategoryDetail{}, err
	}

	accountNames, err := s.loadAccountNames(userID)
	if err != nil {
		return DashboardCategoryDetail{}, err
	}

	analyticsTransactions := filterDashboardAnalyticsTransactions(periodTransactions)
	matchingTransactions := make([]Transaction, 0)
	totalExpenses := 0.0

	for _, transaction := range analyticsTransactions {
		if transaction.Type != Expense {
			continue
		}

		totalExpenses += transaction.Amount

		currentKey, _ := normalizeDashboardCategory(transaction.Category)
		if currentKey == normalizedKey {
			matchingTransactions = append(matchingTransactions, transaction)
		}
	}

	totalSpent := 0.0
	for _, transaction := range matchingTransactions {
		totalSpent += transaction.Amount
	}

	sort.Slice(matchingTransactions, func(left, right int) bool {
		if matchingTransactions[left].CreatedAt.Equal(matchingTransactions[right].CreatedAt) {
			return matchingTransactions[left].ID > matchingTransactions[right].ID
		}
		return matchingTransactions[left].CreatedAt.After(matchingTransactions[right].CreatedAt)
	})

	percentage := 0.0
	if totalExpenses > 0 {
		percentage = (totalSpent / totalExpenses) * 100
	}

	interval := selectTrendInterval(rangeStart, rangeEnd)

	return DashboardCategoryDetail{
		Category: DashboardCategorySummary{
			Key:               normalizedKey,
			Label:             categoryLabel,
			Amount:            totalSpent,
			Percentage:        percentage,
			TransactionsCount: int64(len(matchingTransactions)),
		},
		TrendInterval: interval,
		Trend:         buildExpenseOnlyTrendPoints(matchingTransactions, rangeStart, rangeEnd, interval),
		Transactions:  buildTransactionItems(matchingTransactions, accountNames),
	}, nil
}

func (s *service) loadAccountNames(userID int64) (map[int64]string, error) {
	accountsList, err := s.accountSvc.GetAll(userID)
	if err != nil {
		return nil, err
	}

	accountNames := make(map[int64]string, len(accountsList))
	for _, account := range accountsList {
		accountNames[account.ID] = account.Name
	}

	return accountNames, nil
}

func buildDashboardFilters(base TransactionFilters, dateFrom, dateTo time.Time) TransactionFilters {
	return TransactionFilters{
		AccountID: base.AccountID,
		DateFrom:  &dateFrom,
		DateTo:    &dateTo,
	}
}

func resolveDashboardRange(filters TransactionFilters) (time.Time, time.Time) {
	now := time.Now().UTC()

	if filters.DateFrom == nil && filters.DateTo == nil {
		return startOfDay(now), endOfDay(now)
	}

	if filters.DateFrom == nil && filters.DateTo != nil {
		return startOfDay(*filters.DateTo), endOfDay(*filters.DateTo)
	}

	if filters.DateFrom != nil && filters.DateTo == nil {
		return startOfDay(*filters.DateFrom), endOfDay(*filters.DateFrom)
	}

	return startOfDay(*filters.DateFrom), endOfDay(*filters.DateTo)
}

func filterDashboardAnalyticsTransactions(transactions []Transaction) []Transaction {
	filtered := make([]Transaction, 0, len(transactions))

	for _, transaction := range transactions {
		if isDashboardInternalTransaction(transaction) {
			continue
		}

		filtered = append(filtered, transaction)
	}

	return filtered
}

func isDashboardInternalTransaction(transaction Transaction) bool {
	category := strings.TrimSpace(strings.ToLower(transaction.Category))
	return transaction.TransferID != nil || category == TransferCategory || category == InitialBalanceCategory
}

func filterTransactionsInRange(transactions []Transaction, dateFrom, dateTo time.Time) []Transaction {
	filtered := make([]Transaction, 0, len(transactions))

	for _, transaction := range transactions {
		if transaction.CreatedAt.Before(dateFrom) || transaction.CreatedAt.After(dateTo) {
			continue
		}

		filtered = append(filtered, transaction)
	}

	return filtered
}

func summarizeTransactions(transactions []Transaction) DashboardSummary {
	summary := DashboardSummary{
		TransactionsCount: int64(len(transactions)),
	}

	for _, transaction := range transactions {
		switch transaction.Type {
		case Income:
			summary.IncomeTotal += transaction.Amount
			summary.NetBalance += transaction.Amount
		case Expense:
			summary.ExpenseTotal += transaction.Amount
			summary.NetBalance -= transaction.Amount
		}
	}

	return summary
}

func buildDashboardComparison(currentMonth, previousMonth DashboardSummary) DashboardComparison {
	return DashboardComparison{
		CurrentMonth:     currentMonth,
		PreviousMonth:    previousMonth,
		IncomeChangePct:  calculatePercentageChange(currentMonth.IncomeTotal, previousMonth.IncomeTotal),
		ExpenseChangePct: calculatePercentageChange(currentMonth.ExpenseTotal, previousMonth.ExpenseTotal),
	}
}

func calculatePercentageChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}

	return ((current - previous) / previous) * 100
}

func selectTrendInterval(dateFrom, dateTo time.Time) string {
	days := int(endOfDay(dateTo).Sub(startOfDay(dateFrom)).Hours()/24) + 1

	switch {
	case days <= 14:
		return "daily"
	case days <= 93:
		return "weekly"
	default:
		return "monthly"
	}
}

func buildTrendPoints(transactions []Transaction, dateFrom, dateTo time.Time, interval string) []DashboardTrendPoint {
	points := buildTrendBuckets(dateFrom, dateTo, interval)

	for _, transaction := range transactions {
		for index := range points {
			if transaction.CreatedAt.Before(points[index].StartDate) || transaction.CreatedAt.After(points[index].EndDate) {
				continue
			}

			if transaction.Type == Income {
				points[index].IncomeTotal += transaction.Amount
			} else if transaction.Type == Expense {
				points[index].ExpenseTotal += transaction.Amount
			}
			break
		}
	}

	return points
}

func buildExpenseOnlyTrendPoints(transactions []Transaction, dateFrom, dateTo time.Time, interval string) []DashboardTrendPoint {
	points := buildTrendBuckets(dateFrom, dateTo, interval)

	for _, transaction := range transactions {
		for index := range points {
			if transaction.CreatedAt.Before(points[index].StartDate) || transaction.CreatedAt.After(points[index].EndDate) {
				continue
			}

			points[index].ExpenseTotal += transaction.Amount
			break
		}
	}

	return points
}

func buildTrendBuckets(dateFrom, dateTo time.Time, interval string) []DashboardTrendPoint {
	buckets := make([]DashboardTrendPoint, 0)
	cursor := startOfDay(dateFrom)
	last := endOfDay(dateTo)

	for !cursor.After(last) {
		var bucketStart time.Time
		var bucketEnd time.Time

		switch interval {
		case "weekly":
			bucketStart = cursor
			bucketEnd = endOfDay(minTime(cursor.AddDate(0, 0, 6), last))
			cursor = startOfDay(bucketEnd.AddDate(0, 0, 1))
		case "monthly":
			bucketStart = startOfMonth(cursor)
			bucketEnd = endOfDay(minTime(endOfMonth(cursor), last))
			cursor = startOfDay(bucketEnd.AddDate(0, 0, 1))
		default:
			bucketStart = startOfDay(cursor)
			bucketEnd = endOfDay(cursor)
			cursor = cursor.AddDate(0, 0, 1)
		}

		buckets = append(buckets, DashboardTrendPoint{
			Label:     formatTrendLabel(bucketStart, bucketEnd, interval),
			StartDate: bucketStart,
			EndDate:   bucketEnd,
		})
	}

	return buckets
}

func buildCategoryBreakdown(transactions []Transaction) []DashboardCategorySummary {
	totals := map[string]*DashboardCategorySummary{}
	totalAmount := 0.0

	for _, transaction := range transactions {
		if transaction.Type != Expense {
			continue
		}

		key, label := normalizeDashboardCategory(transaction.Category)
		entry, ok := totals[key]
		if !ok {
			entry = &DashboardCategorySummary{Key: key, Label: label}
			totals[key] = entry
		}

		entry.Amount += transaction.Amount
		entry.TransactionsCount++
		totalAmount += transaction.Amount
	}

	result := make([]DashboardCategorySummary, 0, len(totals))
	for _, category := range totals {
		if totalAmount > 0 {
			category.Percentage = (category.Amount / totalAmount) * 100
		}
		result = append(result, *category)
	}

	sort.Slice(result, func(left, right int) bool {
		if result[left].Amount == result[right].Amount {
			return result[left].Label < result[right].Label
		}
		return result[left].Amount > result[right].Amount
	})

	return result
}

func buildTopExpenses(transactions []Transaction, accountNames map[int64]string) []DashboardTransactionItem {
	expenses := make([]Transaction, 0)

	for _, transaction := range transactions {
		if transaction.Type == Expense {
			expenses = append(expenses, transaction)
		}
	}

	sort.Slice(expenses, func(left, right int) bool {
		if expenses[left].Amount == expenses[right].Amount {
			return expenses[left].CreatedAt.After(expenses[right].CreatedAt)
		}
		return expenses[left].Amount > expenses[right].Amount
	})

	if len(expenses) > 5 {
		expenses = expenses[:5]
	}

	return buildTransactionItems(expenses, accountNames)
}

func buildRecentTransactions(transactions []Transaction, accountNames map[int64]string) []DashboardTransactionItem {
	ordered := append([]Transaction(nil), transactions...)

	sort.Slice(ordered, func(left, right int) bool {
		if ordered[left].CreatedAt.Equal(ordered[right].CreatedAt) {
			return ordered[left].ID > ordered[right].ID
		}
		return ordered[left].CreatedAt.After(ordered[right].CreatedAt)
	})

	if len(ordered) > 6 {
		ordered = ordered[:6]
	}

	return buildTransactionItems(ordered, accountNames)
}

func buildTransactionItems(transactions []Transaction, accountNames map[int64]string) []DashboardTransactionItem {
	items := make([]DashboardTransactionItem, 0, len(transactions))

	for _, transaction := range transactions {
		_, categoryLabel := displayDashboardCategory(transaction)

		items = append(items, DashboardTransactionItem{
			ID:            transaction.ID,
			Title:         transaction.Title,
			Amount:        transaction.Amount,
			Type:          transaction.Type,
			AccountID:     transaction.AccountID,
			AccountName:   accountNames[transaction.AccountID],
			Category:      transaction.Category,
			CategoryLabel: categoryLabel,
			Description:   transaction.Description,
			CreatedAt:     transaction.CreatedAt,
		})
	}

	return items
}

func displayDashboardCategory(transaction Transaction) (string, string) {
	if transaction.TransferID != nil || strings.EqualFold(strings.TrimSpace(transaction.Category), TransferCategory) {
		return TransferCategory, "Transferencia"
	}

	if strings.EqualFold(strings.TrimSpace(transaction.Category), InitialBalanceCategory) {
		return InitialBalanceCategory, "Saldo inicial"
	}

	return normalizeDashboardCategory(transaction.Category)
}

func normalizeDashboardCategory(value string) (string, string) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return DashboardCategoryOther, "Otros"
	}

	for _, category := range dashboardCategories {
		if normalized == category.Key {
			return category.Key, category.Label
		}

		for _, keyword := range category.Keywords {
			if strings.Contains(normalized, keyword) {
				return category.Key, category.Label
			}
		}
	}

	return DashboardCategoryOther, "Otros"
}

func resolveDashboardCategory(value string) (string, string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", "", ErrDashboardCategoryRequired
	}

	for _, category := range dashboardCategories {
		if normalized == category.Key {
			return category.Key, category.Label, nil
		}
	}

	return "", "", ErrInvalidDashboardCategory
}

func comparablePreviousMonthRange(anchor time.Time) (time.Time, time.Time) {
	currentMonthStart := startOfMonth(anchor)
	previousMonthAnchor := currentMonthStart.AddDate(0, -1, 0)
	previousMonthStart := startOfMonth(previousMonthAnchor)
	previousMonthLastDay := endOfMonth(previousMonthAnchor).Day()
	day := anchor.Day()
	if day > previousMonthLastDay {
		day = previousMonthLastDay
	}

	previousMonthEnd := time.Date(
		previousMonthAnchor.Year(),
		previousMonthAnchor.Month(),
		day,
		23,
		59,
		59,
		0,
		anchor.Location(),
	)

	return previousMonthStart, previousMonthEnd
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func endOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, 0, value.Location())
}

func startOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func endOfMonth(value time.Time) time.Time {
	startNextMonth := time.Date(value.Year(), value.Month()+1, 1, 0, 0, 0, 0, value.Location())
	return startNextMonth.Add(-time.Second)
}

func minTime(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}

	return right
}

func formatTrendLabel(start, end time.Time, interval string) string {
	switch interval {
	case "weekly":
		return formatShortDate(start) + " - " + formatShortDate(end)
	case "monthly":
		return formatMonthYear(start)
	default:
		return formatShortDate(start)
	}
}

func formatShortDate(value time.Time) string {
	monthLabels := []string{
		"ene", "feb", "mar", "abr", "may", "jun",
		"jul", "ago", "sep", "oct", "nov", "dic",
	}

	return value.Format("02") + " " + monthLabels[int(value.Month())-1]
}

func formatMonthYear(value time.Time) string {
	monthLabels := []string{
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}

	return monthLabels[int(value.Month())-1] + " " + value.Format("2006")
}
