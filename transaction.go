package main

import "github.com/shopspring/decimal"

type Transaction struct {
	Amount          int64  `json:"amount"`
	Merchant        string `json:"merchant"`
	TransactionTime string `json:"transaction-time"`
}

type TransactionsResponse struct {
	Error           string        `json:"error"`
	TransactionList []Transaction `json:"transactions"`
}

type MonthlyReport struct {
	Spent  decimal.Decimal `json:"spent"`
	Income decimal.Decimal `json:"income"`
}
