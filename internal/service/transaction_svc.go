package service

import (
	"bsnack/internal/model"
	"bsnack/internal/repository"
	customerror "bsnack/pkg/custom_error"
	"context"
	"database/sql"
)

type TransactionService struct {
	DB              *sql.DB
	TransactionRepo repository.TransactionRepository
	CustomerRepo    repository.CustomerRepository
	ProductRepo     repository.ProductRepository
}

func NewTransactionService(
	db *sql.DB,
	transactionRepo repository.TransactionRepository,
	customerRepo repository.CustomerRepository,
	productRepo repository.ProductRepository,
) *TransactionService {
	return &TransactionService{
		DB:              db,
		TransactionRepo: transactionRepo,
		CustomerRepo:    customerRepo,
		ProductRepo:     productRepo,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	if len(req.Items) == 0 {
		return nil, customerror.ErrItemsEmpty
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	customer, isNew, err := s.CustomerRepo.GetOrCreate(ctx, tx, req.CustomerName)
	if err != nil {
		return nil, err
	}

	products := make(map[int64]*model.Product)
	total := 0

	for _, item := range req.Items {
		product, err := s.ProductRepo.GetByIDForUpdate(ctx, tx, item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.Stock < item.Quantity {
			return nil, customerror.ErrStockNotEnough
		}

		products[item.ProductID] = product
		total += product.Price * item.Quantity
	}

	trx := &model.Transaction{
		CustomerID: customer.ID,
		TotalPrice: total,
	}

	if err := s.TransactionRepo.CreateTransaction(ctx, tx, trx); err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		product := products[item.ProductID]

		err := s.TransactionRepo.CreateTransactionItem(ctx, tx, &model.TransactionItem{
			TransactionID: trx.ID,
			ProductID:     product.ID,
			Quantity:      item.Quantity,
			Price:         product.Price,
		})
		if err != nil {
			return nil, err
		}

		if err := s.ProductRepo.DecreaseStock(ctx, tx, product.ID, item.Quantity); err != nil {
			return nil, err
		}
	}

	point := total / 1000
	if err := s.CustomerRepo.AddPoint(ctx, tx, customer.ID, point); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.CreateTransactionResponse{
		TransactionID: trx.ID,
		TotalPrice:    total,
		PointEarned:   point,
		IsNewCustomer: isNew,
	}, nil
}
