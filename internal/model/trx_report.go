package model

import "time"

type TransactionReportResponse struct {
	TotalCustomer      int                    `json:"total_customer"`
	NewCustomer        bool                   `json:"new_customer"`
	TotalIncome        int                    `json:"total_income"`
	BestSeller         string                 `json:"best_seller"`
	TotalProductSold   int                    `json:"total_product_sold"`
	LatestTransactions []TransactionReportRow `json:"latest_transactions"`
}

type TransactionReportRow struct {
	TransactionID string    `json:"transaction_id"`
	CustomerName  string    `json:"customer_name"`
	ProductName   string    `json:"product_name"`
	ProductSize   string    `json:"product_size"`
	ProductFlavor string    `json:"product_flavor"`
	Quantity      int       `json:"quantity"`
	Price         int       `json:"price"`
	TransactionAt time.Time `json:"transaction_at"`
	IsNewCustomer bool      `json:"is_new_customer"`
}
