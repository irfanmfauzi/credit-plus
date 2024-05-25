package repository

import (
	"context"
	"credit-plus/internal/model/request"

	"github.com/jmoiron/sqlx"
)

type installmentRepo struct {
	baseRepo
}

func NewInstallmentRepo(db *sqlx.DB) installmentRepo {
	return installmentRepo{
		baseRepo: baseRepo{db: db},
	}
}

func (i *installmentRepo) BulkInsertInstallment(ctx context.Context, tx TxProvider, req []request.CreateCustomerInstallmentRequest) error {
	insertCustomerInstallmentQuery := "INSERT INTO installments (contract_id, installment_amount, paid_amount, due_date) VALUES (:contract_id, :installment_amount, :paid_amount, :due_date)"
	insertCustomerInstallmentQuery, args, err := sqlx.Named(insertCustomerInstallmentQuery, req)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sqlx.Rebind(sqlx.QUESTION, insertCustomerInstallmentQuery), args...)
	if err != nil {
		return err
	}

	return nil
}
