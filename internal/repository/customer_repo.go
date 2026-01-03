package repository

import (
	"bsnack/internal/model"
	"context"
	"database/sql"
)

type CustomerRepo struct {
	DB *sql.DB
}

type CustomerRepository interface {
	GetOrCreate(ctx context.Context, tx *sql.Tx, name string) (*model.Customer, bool, error)
	AddPoint(ctx context.Context, tx *sql.Tx, customerID int64, point int) error
}

func NewCustomerRepository(db *sql.DB) *CustomerRepo {
	return &CustomerRepo{
		DB: db,
	}
}

func (r *CustomerRepo) GetOrCreate(ctx context.Context, tx *sql.Tx, name string) (*model.Customer, bool, error) {

	var customer model.Customer

	err := tx.QueryRowContext(
		ctx,
		"SELECT id, name, point FROM customers WHERE name=$1",
		name,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Point,
	)

	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(
			ctx,
			"INSERT INTO customers(name) VALUES($1) RETURNING id, name, point",
			name,
		).Scan(
			&customer.ID,
			&customer.Name,
			&customer.Point,
		)

		if err != nil {
			return nil, false, err
		}

		return &customer, true, nil
	}

	if err != nil {
		return nil, false, err
	}

	return &customer, false, nil
}

func (r *CustomerRepo) AddPoint(ctx context.Context, tx *sql.Tx, customerID int64, point int) error {
	_, err := tx.ExecContext(ctx,
		"UPDATE customers SET point = point + $1 WHERE id = $2",
		point, customerID,
	)
	return err
}
