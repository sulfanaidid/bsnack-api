package repository

import (
	"context"
	"database/sql"
	"time"

	"bsnack/internal/model"
)

type ReportRepo struct {
	DB *sql.DB
}

type ReportRepository interface {
	GetTransactionReport(ctx context.Context, startDate time.Time, endDate time.Time) (*model.TransactionReportResponse, error)
}

func NewReportRepository(db *sql.DB) *ReportRepo {
	return &ReportRepo{DB: db}
}

func (r *ReportRepo) GetTransactionReport(ctx context.Context, startDate time.Time, endDate time.Time) (*model.TransactionReportResponse, error) {
	report := &model.TransactionReportResponse{}

	err := r.DB.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT customer_id)
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&report.TotalCustomer)

	if err != nil {
		return nil, err
	}

	var newCustomerCount int
	err = r.DB.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM (
		SELECT customer_id, MIN(created_at) AS first_tx
		FROM transactions
		GROUP BY customer_id
	) t
	WHERE t.first_tx BETWEEN date_trunc('month', $1::timestamp) AND $2
	`, startDate, endDate).Scan(&newCustomerCount)

	if err != nil {
		return nil, err
	}

	report.NewCustomer = newCustomerCount > 0

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_price), 0)
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&report.TotalIncome)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(ti.quantity), 0)
		FROM transaction_items ti
		JOIN transactions t ON t.id = ti.transaction_id
		WHERE t.created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&report.TotalProductSold)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT p.name || ' - ' || p.flavor
		FROM transaction_items ti
		JOIN products p ON p.id = ti.product_id
		JOIN transactions t ON t.id = ti.transaction_id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.name, p.flavor
		ORDER BY SUM(ti.quantity) DESC
		LIMIT 1
	`, startDate, endDate).Scan(&report.BestSeller)

	if err == sql.ErrNoRows {
		report.BestSeller = ""
	} else if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, `
	SELECT
		t.id::text,
		c.name,
		p.name,
		p.size,
		p.flavor,
		ti.quantity,
		ti.price,
		t.created_at,
		(
			SELECT MIN(t2.created_at)
			FROM transactions t2
			WHERE t2.customer_id = t.customer_id
		) >= date_trunc('month', t.created_at) AS is_new_customer
	FROM transactions t
	JOIN customers c ON c.id = t.customer_id
	JOIN transaction_items ti ON ti.transaction_id = t.id
	JOIN products p ON p.id = ti.product_id
	WHERE t.created_at BETWEEN $1 AND $2
	ORDER BY t.created_at DESC
	LIMIT 10
	`, startDate, endDate)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row model.TransactionReportRow
		err := rows.Scan(
			&row.TransactionID,
			&row.CustomerName,
			&row.ProductName,
			&row.ProductSize,
			&row.ProductFlavor,
			&row.Quantity,
			&row.Price,
			&row.TransactionAt,
			&row.IsNewCustomer,
		)

		if err != nil {
			return nil, err
		}

		report.LatestTransactions = append(report.LatestTransactions, row)
	}

	return report, nil
}
