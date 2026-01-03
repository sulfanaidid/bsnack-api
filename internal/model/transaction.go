package model

import "time"

// req model
type CreateTransactionRequest struct {
	CustomerName string                   `json:"customer_name"`
	Items        []TransactionItemRequest `json:"items"`
}

type TransactionItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// res models
type CreateTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	TotalPrice    int    `json:"total_price"`
	PointEarned   int    `json:"point_earned"`
	IsNewCustomer bool   `json:"is_new_customer"`
}

// db models
type Transaction struct {
	ID         string
	CustomerID int64
	TotalPrice int
	CreatedAt  time.Time
}

type TransactionItem struct {
	ID            int64
	TransactionID string
	ProductID     int64
	Quantity      int
	Price         int
}
