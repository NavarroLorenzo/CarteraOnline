package transactions

import (
	"errors"
	"sort"
	"time"
)

var (
	ErrDashboardCategoryRequired = errors.New("category_key es obligatorio")
	ErrInvalidDashboardCategory  = errors.New("category_key inválido")
)

type DashboardSummary struct {
	IncomeTotal       float64 `json:"income_total"`
	ExpenseTotal      float64 `json:"expense_total"`
	NetBalance        float64 `json:"net_balance"`
	TransactionsCount int64   `json:"transactions_count"`
}

type DashboardActivePeriod struct {
	Key      DashboardPeriod `json:"key"`
	Label    string          `json:"label"`
	DateFrom time.Time       `json:"date_from"`
	DateTo   time.Time       `json:"date_to"`
}

type DashboardComparison struct {
	Title            string           `json:"title"`
	CurrentLabel     string           `json:"current_label"`
	PreviousLabel    string           `json:"previous_label"`
	Current          DashboardSummary `json:"current"`
	Previous         DashboardSummary `json:"previous"`
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
	ActivePeriod       DashboardActivePeriod      `json:"active_period"`
	PeriodSummary      DashboardSummary           `json:"period_summary"`
	Comparison         *DashboardComparison       `json:"comparison,omitempty"`
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

func (s *service) GetDashboard(userID int64, filters TransactionFilters) (DashboardAnalytics, error) {
	activePeriodKey, rangeStart, rangeEnd := resolveDashboardRange(filters)

	periodTransactions, err := s.getTransactions(userID, buildDashboardFilters(filters, rangeStart, rangeEnd))
	if err != nil {
		return DashboardAnalytics{}, err
	}

	accountNames, err := s.loadAccountNames(userID)
	if err != nil {
		return DashboardAnalytics{}, err
	}

	analyticsTransactions := filterDashboardAnalyticsTransactions(periodTransactions)
	comparison := (*DashboardComparison)(nil)
	comparisonStart, comparisonEnd, hasComparison := resolveDashboardComparisonRange(activePeriodKey, rangeStart, rangeEnd)
	if hasComparison {
		previousPeriodTransactions, err := s.getTransactions(
			userID,
			buildDashboardFilters(filters, comparisonStart, comparisonEnd),
		)
		if err != nil {
			return DashboardAnalytics{}, err
		}

		currentSummary := summarizeTransactions(analyticsTransactions)
		previousSummary := summarizeTransactions(filterDashboardAnalyticsTransactions(previousPeriodTransactions))
		comparison = buildDashboardComparison(
			activePeriodKey,
			currentSummary,
			previousSummary,
		)
	}

	recentTransactions, err := s.getTransactions(userID, buildRecentDashboardFilters(filters))
	if err != nil {
		return DashboardAnalytics{}, err
	}
	trendInterval := selectTrendInterval(rangeStart, rangeEnd)

	return DashboardAnalytics{
		ActivePeriod:       buildDashboardActivePeriod(activePeriodKey, rangeStart, rangeEnd),
		PeriodSummary:      summarizeTransactions(analyticsTransactions),
		Comparison:         comparison,
		TrendInterval:      trendInterval,
		Trend:              buildTrendPoints(analyticsTransactions, rangeStart, rangeEnd, trendInterval),
		ExpenseCategories:  buildCategoryBreakdown(analyticsTransactions),
		TopExpenses:        buildTopExpenses(analyticsTransactions, accountNames),
		RecentTransactions: buildRecentTransactions(recentTransactions, accountNames),
	}, nil
}

func (s *service) GetDashboardCategoryDetail(userID int64, filters TransactionFilters, categoryKey string) (DashboardCategoryDetail, error) {
	normalizedKey, categoryLabel, err := resolveDashboardCategory(categoryKey)
	if err != nil {
		return DashboardCategoryDetail{}, err
	}

	_, rangeStart, rangeEnd := resolveDashboardRange(filters)
	periodTransactions, err := s.getTransactions(userID, buildDashboardFilters(filters, rangeStart, rangeEnd))
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

		if transaction.Category == normalizedKey {
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
		Type:      base.Type,
		Category:  base.Category,
		Period:    base.Period,
		DateFrom:  &dateFrom,
		DateTo:    &dateTo,
	}
}

func buildRecentDashboardFilters(base TransactionFilters) TransactionFilters {
	return TransactionFilters{
		AccountID: base.AccountID,
		Type:      base.Type,
		Category:  base.Category,
	}
}

func resolveDashboardRange(filters TransactionFilters) (DashboardPeriod, time.Time, time.Time) {
	now := time.Now().UTC()
	anchor := endOfDay(now)

	if filters.DateTo != nil {
		anchor = endOfDay(*filters.DateTo)
	} else if filters.DateFrom != nil {
		anchor = endOfDay(*filters.DateFrom)
	}

	if filters.Period != nil {
		switch *filters.Period {
		case DashboardPeriodDay:
			return DashboardPeriodDay, startOfDay(anchor), endOfDay(anchor)
		case DashboardPeriodWeek:
			return DashboardPeriodWeek, startOfWeek(anchor), endOfDay(anchor)
		case DashboardPeriodYear:
			return DashboardPeriodYear, startOfYear(anchor), endOfDay(anchor)
		case DashboardPeriodCustom:
			if filters.DateFrom != nil && filters.DateTo != nil {
				return DashboardPeriodCustom, startOfDay(*filters.DateFrom), endOfDay(*filters.DateTo)
			}
		case DashboardPeriodMonth:
			return DashboardPeriodMonth, startOfMonth(anchor), endOfDay(anchor)
		}
	}

	if filters.DateFrom == nil && filters.DateTo == nil {
		return DashboardPeriodMonth, startOfMonth(anchor), endOfDay(anchor)
	}

	if filters.DateFrom == nil && filters.DateTo != nil {
		return DashboardPeriodDay, startOfDay(*filters.DateTo), endOfDay(*filters.DateTo)
	}

	if filters.DateFrom != nil && filters.DateTo == nil {
		return DashboardPeriodDay, startOfDay(*filters.DateFrom), endOfDay(*filters.DateFrom)
	}

	return DashboardPeriodCustom, startOfDay(*filters.DateFrom), endOfDay(*filters.DateTo)
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
	return transaction.TransferID != nil || transaction.Category == CategoryTransfer
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

func buildDashboardActivePeriod(period DashboardPeriod, dateFrom, dateTo time.Time) DashboardActivePeriod {
	return DashboardActivePeriod{
		Key:      period,
		Label:    dashboardPeriodLabel(period),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
}

func buildDashboardComparison(
	period DashboardPeriod,
	current DashboardSummary,
	previous DashboardSummary,
) *DashboardComparison {
	title, currentLabel, previousLabel := dashboardComparisonLabels(period)

	return &DashboardComparison{
		Title:            title,
		CurrentLabel:     currentLabel,
		PreviousLabel:    previousLabel,
		Current:          current,
		Previous:         previous,
		IncomeChangePct:  calculatePercentageChange(current.IncomeTotal, previous.IncomeTotal),
		ExpenseChangePct: calculatePercentageChange(current.ExpenseTotal, previous.ExpenseTotal),
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
	case days <= 1:
		return "hourly"
	case days <= 14:
		return "daily"
	case days <= 31:
		return "weekly"
	default:
		return "monthly"
	}
}

func buildTrendPoints(transactions []Transaction, dateFrom, dateTo time.Time, interval string) []DashboardTrendPoint {
	if interval == "hourly" {
		return buildHourlyTrendPoints(transactions, dateFrom, dateTo)
	}

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
	if interval == "hourly" {
		hourlyPoints := buildHourlyTrendPoints(transactions, dateFrom, dateTo)

		for index := range hourlyPoints {
			hourlyPoints[index].IncomeTotal = 0
		}

		return hourlyPoints
	}

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

		entry, ok := totals[transaction.Category]
		if !ok {
			entry = &DashboardCategorySummary{
				Key:   transaction.Category,
				Label: resolveTransactionCategoryLabel(transaction.Category),
			}
			totals[transaction.Category] = entry
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

	if len(ordered) > 5 {
		ordered = ordered[:5]
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

func buildHourlyTrendPoints(transactions []Transaction, dateFrom, dateTo time.Time) []DashboardTrendPoint {
	type pointTotals struct {
		income  float64
		expense float64
	}

	groupedTotals := make(map[time.Time]pointTotals)

	for _, transaction := range transactions {
		bucketStart := transaction.CreatedAt.Truncate(time.Hour)
		if bucketStart.Before(startOfDay(dateFrom)) || bucketStart.After(endOfDay(dateTo)) {
			continue
		}

		totals := groupedTotals[bucketStart]
		if transaction.Type == Income {
			totals.income += transaction.Amount
		} else if transaction.Type == Expense {
			totals.expense += transaction.Amount
		}

		groupedTotals[bucketStart] = totals
	}

	if len(groupedTotals) <= 1 {
		return buildSingleDayTrendPoint(transactions, dateFrom)
	}

	keys := make([]time.Time, 0, len(groupedTotals))
	for bucketStart := range groupedTotals {
		keys = append(keys, bucketStart)
	}

	sort.Slice(keys, func(left, right int) bool {
		return keys[left].Before(keys[right])
	})

	points := make([]DashboardTrendPoint, 0, len(keys))
	for _, bucketStart := range keys {
		totals := groupedTotals[bucketStart]
		points = append(points, DashboardTrendPoint{
			Label:        bucketStart.Format("15h"),
			StartDate:    bucketStart,
			EndDate:      bucketStart.Add(time.Hour - time.Second),
			IncomeTotal:  totals.income,
			ExpenseTotal: totals.expense,
		})
	}

	return points
}

func buildSingleDayTrendPoint(transactions []Transaction, dateFrom time.Time) []DashboardTrendPoint {
	point := DashboardTrendPoint{
		Label:     "Hoy",
		StartDate: startOfDay(dateFrom),
		EndDate:   endOfDay(dateFrom),
	}

	for _, transaction := range transactions {
		if transaction.Type == Income {
			point.IncomeTotal += transaction.Amount
		} else if transaction.Type == Expense {
			point.ExpenseTotal += transaction.Amount
		}
	}

	return []DashboardTrendPoint{point}
}

func displayDashboardCategory(transaction Transaction) (string, string) {
	if transaction.TransferID != nil || transaction.Category == CategoryTransfer {
		return CategoryTransfer, resolveTransactionCategoryLabel(CategoryTransfer)
	}

	if isInitialBalanceTransaction(transaction) {
		return CategoryOther, "Saldo inicial"
	}

	return transaction.Category, resolveTransactionCategoryLabel(transaction.Category)
}

func resolveDashboardCategory(value string) (string, string, error) {
	if normalizeCategoryToken(value) == "" {
		return "", "", ErrDashboardCategoryRequired
	}

	key, found := resolveTransactionCategoryKey(value)
	if !found {
		return "", "", ErrInvalidDashboardCategory
	}

	return key, resolveTransactionCategoryLabel(key), nil
}

func resolveDashboardComparisonRange(period DashboardPeriod, currentStart, currentEnd time.Time) (time.Time, time.Time, bool) {
	switch period {
	case DashboardPeriodDay:
		previousDay := currentStart.AddDate(0, 0, -1)
		return startOfDay(previousDay), endOfDay(previousDay), true
	case DashboardPeriodWeek:
		previousEnd := endOfDay(currentStart.AddDate(0, 0, -1))
		return startOfDay(previousEnd.AddDate(0, 0, -(inclusiveDaysBetween(currentStart, currentEnd) - 1))), previousEnd, true
	case DashboardPeriodMonth:
		previousStart, previousEnd := comparablePreviousMonthRange(currentEnd)
		return previousStart, previousEnd, true
	case DashboardPeriodYear:
		previousStart, previousEnd := comparablePreviousYearRange(currentEnd)
		return previousStart, previousEnd, true
	case DashboardPeriodCustom:
		previousEnd := endOfDay(currentStart.AddDate(0, 0, -1))
		return startOfDay(previousEnd.AddDate(0, 0, -(inclusiveDaysBetween(currentStart, currentEnd) - 1))), previousEnd, true
	default:
		return time.Time{}, time.Time{}, false
	}
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

func comparablePreviousYearRange(anchor time.Time) (time.Time, time.Time) {
	previousYearAnchor := anchor.AddDate(-1, 0, 0)
	return startOfYear(previousYearAnchor), endOfDay(previousYearAnchor)
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func endOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, 0, value.Location())
}

func startOfWeek(value time.Time) time.Time {
	weekday := int(value.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	return startOfDay(value.AddDate(0, 0, -(weekday - 1)))
}

func startOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func startOfYear(value time.Time) time.Time {
	return time.Date(value.Year(), time.January, 1, 0, 0, 0, 0, value.Location())
}

func endOfMonth(value time.Time) time.Time {
	startNextMonth := time.Date(value.Year(), value.Month()+1, 1, 0, 0, 0, 0, value.Location())
	return startNextMonth.Add(-time.Second)
}

func inclusiveDaysBetween(start, end time.Time) int {
	return int(endOfDay(end).Sub(startOfDay(start)).Hours()/24) + 1
}

func minTime(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}

	return right
}

func formatTrendLabel(start, end time.Time, interval string) string {
	switch interval {
	case "hourly":
		return start.Format("15h")
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

func dashboardPeriodLabel(period DashboardPeriod) string {
	switch period {
	case DashboardPeriodDay:
		return "Día"
	case DashboardPeriodWeek:
		return "Semana"
	case DashboardPeriodYear:
		return "Año"
	case DashboardPeriodCustom:
		return "Rango personalizado"
	case DashboardPeriodMonth:
		fallthrough
	default:
		return "Mes"
	}
}

func dashboardComparisonLabels(period DashboardPeriod) (string, string, string) {
	switch period {
	case DashboardPeriodDay:
		return "Día actual vs anterior", "Día actual", "Día anterior"
	case DashboardPeriodWeek:
		return "Semana actual vs anterior", "Semana actual", "Semana anterior"
	case DashboardPeriodYear:
		return "Año actual vs anterior", "Año actual", "Año anterior"
	case DashboardPeriodCustom:
		return "Período actual vs anterior", "Período actual", "Período anterior"
	case DashboardPeriodMonth:
		fallthrough
	default:
		return "Mes actual vs anterior", "Mes actual", "Mes anterior"
	}
}
