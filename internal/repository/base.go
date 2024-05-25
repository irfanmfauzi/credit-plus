package repository

import "github.com/jmoiron/sqlx"

type baseRepo struct {
	db *sqlx.DB
}

func (b *baseRepo) DB(tx TxProvider) QueryProvider {
	if tx != nil {
		return tx
	}
	return b.db
}
