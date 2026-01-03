package repository

import (
	"bsnack/internal/model"
	"context"
	"database/sql"
)

type TransactionRepo struct {
	DB *sql.DB
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, t *model.Transaction) error
	CreateTransactionItem(ctx context.Context, tx *sql.Tx, txItem *model.TransactionItem) error
}

func NewTransactionRepository(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{
		DB: db,
	}
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, tx *sql.Tx, t *model.Transaction) error {
	return tx.QueryRowContext(ctx,
		"INSERT INTO transactions(customer_id, total_price) VALUES($1,$2) RETURNING id",
		t.CustomerID, t.TotalPrice,
	).Scan(&t.ID)
}

func (r *TransactionRepo) CreateTransactionItem(ctx context.Context, tx *sql.Tx, item *model.TransactionItem) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO transaction_items
		(transaction_id, product_id, quantity, price)
		VALUES ($1,$2,$3,$4)`,
		item.TransactionID, item.ProductID, item.Quantity, item.Price,
	)
	return err
}
