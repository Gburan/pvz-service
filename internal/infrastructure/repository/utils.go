package repository

import "errors"

var (
	ErrBuildQuery         = errors.New("failed to build SQL query")
	ErrExecuteQuery       = errors.New("failed to execute query")
	ErrScanResult         = errors.New("failed to scan result")
	ErrReceptionNotFound  = errors.New("reception not found")
	ErrReceptionsNotFound = errors.New("receptions not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductsNotFound   = errors.New("products not found")
	ErrPVZNotFound        = errors.New("pvz not found")
	ErrUserNotFound       = errors.New("user not found")
)
