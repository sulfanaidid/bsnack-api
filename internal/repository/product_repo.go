package repository

import (
	"bsnack/internal/model"
	customerror "bsnack/pkg/custom_error"
	"context"
	"database/sql"
)

type ProductsRepo struct {
	DB *sql.DB
}

type ProductRepository interface {
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.Product, error)
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*model.Product, error)
	DecreaseStock(ctx context.Context, tx *sql.Tx, productID int64, qty int) error
}

func NewProductRepository(db *sql.DB) *ProductsRepo {
	return &ProductsRepo{
		DB: db,
	}
}

func (r *ProductsRepo) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.Product, error) {
	var product model.Product

	err := tx.QueryRowContext(
		ctx,
		`SELECT id, price, stock
		 FROM products
		 WHERE id = $1 AND is_active = true`,
		id,
	).Scan(
		&product.ID,
		&product.Price,
		&product.Stock,
	)

	if err == sql.ErrNoRows {
		return nil, customerror.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductsRepo) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*model.Product, error) {
	var product model.Product

	err := tx.QueryRowContext(
		ctx,
		`SELECT id, price, stock
		 FROM products
		 WHERE id = $1 AND is_active = true
		 FOR UPDATE`,
		id,
	).Scan(
		&product.ID,
		&product.Price,
		&product.Stock,
	)

	if err == sql.ErrNoRows {
		return nil, customerror.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductsRepo) DecreaseStock(ctx context.Context, tx *sql.Tx, productID int64, qty int) error {
	_, err := tx.ExecContext(ctx,
		"UPDATE products SET stock = stock - $1 WHERE id = $2",
		qty, productID,
	)
	return err
}
