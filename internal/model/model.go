package model

import "context"

type Repository interface {
}

// TxFunc represents function to run in an SQL-transaction.
type TxFunc func(ctx context.Context, tx Repository) (err error)

type DB interface {
	// Repository access without transactions
	Repository
	// InTx runs given function in transaction
	InTx(context.Context, TxFunc) error
}
