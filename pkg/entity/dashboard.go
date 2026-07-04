package entity

import "time"

type DashboardFilterDto struct {
	CustomerID      string    `query:"customerId"`
	StartDate     time.Time `query:"startDate"`
	EndDate      time.Time `query:"endDate"`
	TimeRange    string    `query:"timeRange"` // 7d, 30d, 3m, 1y
}

type CashFlowData struct {
	Date     string  `json:"date"`
	Income   float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type ExpenseByCategory struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color"`
	Total       float64 `json:"total"`
	Percentage  float64 `json:"percentage"`
}

type DailySales struct {
	Date            string  `json:"date"`
	DayOfWeek       int     `json:"dayOfWeek"`
	DailyRevenue    float64 `json:"dailyRevenue"`
	TransactionCount int    `json:"transactionCount"`
	Average       float64 `json:"average"`
}

type ProfitLossData struct {
	Label        string  `json:"label"`
	Value       float64 `json:"value"`
	Type        string  `json:"type"` // INCOME, EXPENSE, NET
}

type DashboardSummaryDto struct {
	TotalIncome     float64           `json:"totalIncome"`
	TotalExpense  float64           `json:"totalExpense"`
	Balance      float64           `json:"balance"`
	CashFlow     []CashFlowData     `json:"cashFlow"`
	ExpensesByCategory []ExpenseByCategory `json:"expensesByCategory"`
	DailySales   []DailySales     `json:"dailySales"`
	ProfitLoss  []ProfitLossData `json:"profitLoss"`
}