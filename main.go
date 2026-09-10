package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Expense struct {
	ID          int
	Icon        string
	Description string
	Category    string
	Date        string
	Amount      float64
}

type CategoryData struct {
	Name    string
	Amount  float64
	Class   string
	Icon    string
	Percent float64
}

type MonthData struct {
	Label  string
	Amount float64
	Height float64
}

type DashboardData struct {
	Expenses       []Expense
	TotalSpent     float64
	TotalExpenses  int
	AverageExpense float64
	MonthlyBudget  float64
	MonthlySpent   float64
	BudgetPercent  float64
	Categories     []CategoryData
	Months         []MonthData
}

type ExpenseRequest struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}

type BudgetRequest struct {
	Budget float64 `json:"budget"`
}

var db *sql.DB

var categoryIcons = map[string]string{
	"Food":          "🍔",
	"Transport":     "🚗",
	"Shopping":      "🛍️",
	"Entertainment": "🎬",
	"Bills":         "💡",
	"Health":        "❤️",
	"Education":     "📚",
	"Other":         "📦",
}

var categoryClasses = map[string]string{
	"Food":          "food",
	"Transport":     "transport",
	"Shopping":      "shopping",
	"Entertainment": "entertainment",
	"Bills":         "bills",
	"Health":        "health",
	"Education":     "education",
	"Other":         "other",
}

func main() {

	var err error

	db, err = sql.Open("sqlite", "expenses.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	createTables()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/add-expense", addExpenseHandler)
	http.HandleFunc("/update-expense", updateExpenseHandler)
	http.HandleFunc("/delete-expense", deleteExpenseHandler)
	http.HandleFunc("/set-budget", setBudgetHandler)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	fmt.Println("======================================")
	fmt.Println("          ExpenseFlow v1.1")
	fmt.Println("      http://localhost:8080")
	fmt.Println("======================================")

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}

func createTables() {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT NOT NULL,
			category TEXT NOT NULL,
			date TEXT NOT NULL,
			amount REAL NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY,
			monthly_budget REAL NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	var count int

	err = db.QueryRow(
		`SELECT COUNT(*) FROM settings`,
	).Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {

		_, err = db.Exec(`
			INSERT INTO settings
			(id, monthly_budget)
			VALUES (1, 30000)
		`)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func dashboardHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	expenses := getExpenses()

	budget := getBudget()

	totalSpent := 0.0

	for _, expense := range expenses {
		totalSpent += expense.Amount
	}

	totalExpenses := len(expenses)

	average := 0.0

	if totalExpenses > 0 {
		average = totalSpent / float64(totalExpenses)
	}

	currentMonth := time.Now().Format("2006-01")

	monthlySpent := 0.0

	for _, expense := range expenses {

		if strings.HasPrefix(
			expense.Date,
			currentMonth,
		) {
			monthlySpent += expense.Amount
		}
	}

	budgetPercent := 0.0

	if budget > 0 {
		budgetPercent = (monthlySpent / budget) * 100
	}

	if budgetPercent > 100 {
		budgetPercent = 100
	}

	data := DashboardData{
		Expenses:       expenses,
		TotalSpent:     totalSpent,
		TotalExpenses:  totalExpenses,
		AverageExpense: average,
		MonthlyBudget:  budget,
		MonthlySpent:   monthlySpent,
		BudgetPercent:  budgetPercent,
		Categories:     buildCategories(expenses),
		Months:         buildMonths(expenses),
	}

	tmpl, err := template.ParseFiles(
		"templates/index.html",
	)

	if err != nil {

		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	err = tmpl.Execute(w, data)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func getExpenses() []Expense {

	rows, err := db.Query(`
		SELECT
			id,
			description,
			category,
			date,
			amount
		FROM expenses
		ORDER BY date DESC, id DESC
	`)

	if err != nil {

		log.Println(err)

		return []Expense{}
	}

	defer rows.Close()

	var expenses []Expense

	for rows.Next() {

		var e Expense

		err := rows.Scan(
			&e.ID,
			&e.Description,
			&e.Category,
			&e.Date,
			&e.Amount,
		)

		if err != nil {

			log.Println(err)

			continue
		}

		e.Icon = categoryIcons[e.Category]

		if e.Icon == "" {
			e.Icon = "📦"
		}

		expenses = append(expenses, e)
	}

	return expenses
}

func getBudget() float64 {

	var budget float64

	err := db.QueryRow(`
		SELECT monthly_budget
		FROM settings
		WHERE id = 1
	`).Scan(&budget)

	if err != nil {
		return 30000
	}

	return budget
}

func addExpenseHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	var request ExpenseRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {

		http.Error(
			w,
			"Invalid request",
			http.StatusBadRequest,
		)

		return
	}

	request.Description = strings.TrimSpace(
		request.Description,
	)

	request.Category = strings.TrimSpace(
		request.Category,
	)

	request.Date = strings.TrimSpace(
		request.Date,
	)

	if request.Description == "" ||
		request.Amount <= 0 ||
		request.Category == "" ||
		request.Date == "" {

		http.Error(
			w,
			"Invalid expense data",
			http.StatusBadRequest,
		)

		return
	}

	icon := categoryIcons[request.Category]

	if icon == "" {
		icon = "📦"
	}

	_, err = db.Exec(`
    INSERT INTO expenses
    (icon, description, category, date, amount)
    VALUES (?, ?, ?, ?, ?)
`,
		icon,
		request.Description,
		request.Category,
		request.Date,
		request.Amount,
	)

	if err != nil {

		log.Println(err)

		http.Error(
			w,
			"Could not save expense",
			http.StatusInternalServerError,
		)

		return
	}

	jsonResponse(
		w,
		"Expense added successfully",
	)
}

func updateExpenseHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPut {

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)

	if err != nil || id <= 0 {

		http.Error(
			w,
			"Invalid expense ID",
			http.StatusBadRequest,
		)

		return
	}

	var request ExpenseRequest

	err = json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {

		http.Error(
			w,
			"Invalid request",
			http.StatusBadRequest,
		)

		return
	}

	request.Description = strings.TrimSpace(
		request.Description,
	)

	request.Category = strings.TrimSpace(
		request.Category,
	)

	request.Date = strings.TrimSpace(
		request.Date,
	)

	if request.Description == "" ||
		request.Amount <= 0 ||
		request.Category == "" ||
		request.Date == "" {

		http.Error(
			w,
			"Invalid expense data",
			http.StatusBadRequest,
		)

		return
	}

	result, err := db.Exec(`
		UPDATE expenses
		SET
			description = ?,
			category = ?,
			date = ?,
			amount = ?
		WHERE id = ?
	`,
		request.Description,
		request.Category,
		request.Date,
		request.Amount,
		id,
	)

	if err != nil {

		log.Println(err)

		http.Error(
			w,
			"Could not update expense",
			http.StatusInternalServerError,
		)

		return
	}

	rows, err := result.RowsAffected()

	if err != nil {

		http.Error(
			w,
			"Could not check update",
			http.StatusInternalServerError,
		)

		return
	}

	if rows == 0 {

		http.Error(
			w,
			"Expense not found",
			http.StatusNotFound,
		)

		return
	}

	jsonResponse(
		w,
		"Expense updated successfully",
	)
}

func deleteExpenseHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodDelete {

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)

	if err != nil || id <= 0 {

		http.Error(
			w,
			"Invalid expense ID",
			http.StatusBadRequest,
		)

		return
	}

	result, err := db.Exec(
		`DELETE FROM expenses WHERE id = ?`,
		id,
	)

	if err != nil {

		log.Println(err)

		http.Error(
			w,
			"Could not delete expense",
			http.StatusInternalServerError,
		)

		return
	}

	rows, err := result.RowsAffected()

	if err != nil {

		http.Error(
			w,
			"Could not check deletion",
			http.StatusInternalServerError,
		)

		return
	}

	if rows == 0 {

		http.Error(
			w,
			"Expense not found",
			http.StatusNotFound,
		)

		return
	}

	jsonResponse(
		w,
		"Expense deleted successfully",
	)
}

func setBudgetHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	var request BudgetRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {

		http.Error(
			w,
			"Invalid request",
			http.StatusBadRequest,
		)

		return
	}

	if request.Budget <= 0 {

		http.Error(
			w,
			"Budget must be greater than zero",
			http.StatusBadRequest,
		)

		return
	}

	_, err = db.Exec(`
		UPDATE settings
		SET monthly_budget = ?
		WHERE id = 1
	`,
		request.Budget,
	)

	if err != nil {

		log.Println(err)

		http.Error(
			w,
			"Could not update budget",
			http.StatusInternalServerError,
		)

		return
	}

	jsonResponse(
		w,
		"Monthly budget updated",
	)
}

func jsonResponse(
	w http.ResponseWriter,
	message string,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": message,
		},
	)
}

func buildCategories(
	expenses []Expense,
) []CategoryData {

	categoryTotals := make(
		map[string]float64,
	)

	for _, expense := range expenses {

		categoryTotals[expense.Category] +=
			expense.Amount
	}

	total := 0.0

	for _, amount := range categoryTotals {
		total += amount
	}

	var categories []CategoryData

	order := []string{
		"Food",
		"Transport",
		"Shopping",
		"Entertainment",
		"Bills",
		"Health",
		"Education",
		"Other",
	}

	for _, category := range order {

		amount := categoryTotals[category]

		if amount == 0 {
			continue
		}

		percent := 0.0

		if total > 0 {
			percent = (amount / total) * 100
		}

		categories = append(
			categories,
			CategoryData{
				Name:    category,
				Amount:  amount,
				Class:   categoryClasses[category],
				Icon:    categoryIcons[category],
				Percent: percent,
			},
		)
	}

	return categories
}

func buildMonths(
	expenses []Expense,
) []MonthData {

	now := time.Now()

	var months []MonthData

	maxAmount := 0.0

	type monthAmount struct {
		label  string
		amount float64
	}

	var values []monthAmount

	for i := 5; i >= 0; i-- {

		date := now.AddDate(
			0,
			-i,
			0,
		)

		monthKey := date.Format(
			"2006-01",
		)

		label := date.Format("Jan")

		amount := 0.0

		for _, expense := range expenses {

			if strings.HasPrefix(
				expense.Date,
				monthKey,
			) {

				amount += expense.Amount
			}
		}

		if amount > maxAmount {
			maxAmount = amount
		}

		values = append(
			values,
			monthAmount{
				label:  label,
				amount: amount,
			},
		)
	}

	for _, value := range values {

		height := 5.0

		if maxAmount > 0 {

			height =
				(value.amount / maxAmount) * 100

			if height < 5 {
				height = 5
			}
		}

		months = append(
			months,
			MonthData{
				Label:  value.label,
				Amount: value.amount,
				Height: height,
			},
		)
	}

	return months
}
