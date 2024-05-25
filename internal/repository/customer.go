package repository

import (
	"context"
	"credit-plus/internal/model/request"

	"github.com/jmoiron/sqlx"
)

type customerRepo struct {
	baseRepo
}

func NewCustomerRepo(db *sqlx.DB) customerRepo {
	return customerRepo{
		baseRepo: baseRepo{db: db},
	}
}

func (c *customerRepo) InsertCustomer(ctx context.Context, tx TxProvider, req request.CreateCustomerRequest) error {
	query := "INSERT INTO customers (nik, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp, foto_selfie, limit_transaction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, req.Nik, req.FullName, req.LegalName, req.TempatLahir, req.TanggalLahir, req.Gaji, req.FotoKTP, req.FotoSelfie, req.Limit)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) CheckLimitTransaction(ctx context.Context, tx TxProvider, otrPrice int, customerId int64) (bool, error) {
	checkLimitQuery := "SELECT limit_transaction > ? FROM customers where id = ?"
	isAllowed := 0

	err := c.DB(tx).GetContext(ctx, &isAllowed, checkLimitQuery, otrPrice, customerId)
	if err != nil {
		return false, err
	}
	return isAllowed == 1, nil
}

func (c *customerRepo) UpdateLimitTransaction(ctx context.Context, tx TxProvider, otrPrice int, customerId int64) error {
	updateLimit := "UPDATE customers set limit_transaction = limit_transaction - ? where id = ?"

	_, err := c.DB(tx).ExecContext(ctx, updateLimit, otrPrice, customerId)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) InsertContract(ctx context.Context, tx TxProvider, req request.CreateContactRequest) (int64, error) {
	var contractId int64

	insertContractQuery := "INSERT INTO contracts (customer_id, no_kontrak, otr_price, tenor, bunga, nama_aset, fee) VALUES (?, ?, ?, ?, ?, ?, ?)"

	result, err := c.DB(tx).ExecContext(ctx, insertContractQuery, req.CustomerID, req.NoKontrak, req.OtrPrice, req.Tenor, req.Bunga, req.NamaAset, req.Fee)
	if err != nil {
		return 0, err
	}

	contractId, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return contractId, nil
}
