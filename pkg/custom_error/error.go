package customerror

import "errors"

var (
	ErrStockNotEnough   = errors.New("stock not enough")
	ErrProductNotFound  = errors.New("product not found")
	ErrItemsEmpty       = errors.New("items required")
	ErrInvalidDateRange = errors.New("invalid date range")
)
