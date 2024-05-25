package repository

import (
	"context"
	"credit-plus/internal/model/request"
	"database/sql"
)

type TxProvider interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type QueryProvider interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type CustomerRepository interface {
	CheckLimitTransaction(ctx context.Context, tx TxProvider, otrPrice int, customerId int64) (bool, error)
	UpdateLimitTransaction(ctx context.Context, tx TxProvider, otrPrice int, customerId int64) error
	InsertContract(ctx context.Context, tx TxProvider, req request.CreateContactRequest) (int64, error)
	InsertCustomer(ctx context.Context, tx TxProvider, req request.CreateCustomerRequest) error
}

type InstallmentRepository interface {
	BulkInsertInstallment(ctx context.Context, tx TxProvider, req []request.CreateCustomerInstallmentRequest) error
}
