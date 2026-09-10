# 💰 ExpenseFlow

<p align="center">
  <strong>A modern web-based expense tracking application built with Go and SQLite.</strong>
</p>

<p align="center">
  Track expenses, manage your monthly budget, analyze spending, and keep your finances organized — all from a clean and intuitive dashboard.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/SQLite-Database-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite">
  <img src="https://img.shields.io/badge/HTML5-Frontend-E34F26?style=for-the-badge&logo=html5&logoColor=white" alt="HTML5">
  <img src="https://img.shields.io/badge/CSS3-Styling-1572B6?style=for-the-badge&logo=css3&logoColor=white" alt="CSS3">
  <img src="https://img.shields.io/badge/JavaScript-Frontend-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black" alt="JavaScript">
</p>

---

## 📌 About The Project

**ExpenseFlow** is a full-stack expense management web application designed to make personal expense tracking simple and organized.

The application provides a dashboard where users can record daily expenses, categorize spending, monitor their monthly budget, and view spending insights.

The backend is developed using **Go**, while **SQLite** is used for persistent data storage. The frontend combines **HTML, CSS, and JavaScript** to provide an interactive web interface.

---

## ✨ Features

### 📊 Dashboard

- Total amount spent
- Total number of expenses
- Average expense
- Monthly budget
- Monthly spending
- Budget usage percentage
- Recent expenses
- Spending overview for the last 6 months

### ➕ Expense Management

- Add new expenses
- Edit existing expenses
- Delete expenses
- Select expense categories
- Choose expense dates
- Store expenses permanently in SQLite

### 🗂️ Categories

ExpenseFlow supports multiple expense categories:

- 🍔 Food
- 🚗 Transport
- 🛍️ Shopping
- 🎬 Entertainment
- 💡 Bills
- ❤️ Health
- 📚 Education
- 📦 Other

### 💰 Budget Management

- Set a monthly budget
- Track monthly spending
- View budget utilization
- Visual budget progress indicator
- Budget status notification

### 🔎 Search & Filtering

- Search expenses by description
- Filter expenses by category
- Quickly find specific transactions

### 📈 Spending Analytics

- Category-wise spending breakdown
- Percentage distribution of expenses
- Six-month spending overview
- Visual spending indicators

### 💾 Persistent Storage

All expenses and budget information are stored in a local **SQLite database**, allowing data to remain available after restarting the application.

---

## 🖥️ Application Preview

<p align="center">
  <img src="screenshots/dashboard.png" width="900" alt="ExpenseFlow Dashboard">
</p>

<p align="center">
  <img src="screenshots/add-expense.png" width="700" alt="ExpenseFlow Add Expense">
</p>

> 📌 Add your actual screenshots inside a `screenshots` folder and update the filenames above if required.

---

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go** | Backend & HTTP server |
| **SQLite** | Database |
| **HTML5** | Application structure |
| **CSS3** | UI styling |
| **JavaScript** | Frontend interactions |
| **Go Templates** | Dynamic HTML rendering |
| **Font Awesome** | Icons |
| **Google Fonts** | Typography |

---

## 🏗️ Project Architecture

```text
ExpenseFlow
│
├── main.go
│
├── go.mod
├── go.sum
│
├── expenses.db
│
├── templates
│   └── index.html
│
└── static
    ├── style.css
    └── script.js
