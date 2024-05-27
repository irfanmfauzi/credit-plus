package repository_test

import (
	"context"
	"credit-plus/internal/model/request"
	"credit-plus/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_customerRepo_InsertCustomer(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		req request.CreateCustomerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "best_case",
			args: args{
				ctx: ctx,
				req: request.CreateCustomerRequest{
					Nik:          "NIK",
					FullName:     "FullName",
					LegalName:    "LegalName",
					TempatLahir:  "Bandung",
					TanggalLahir: "24 Juli 1996",
					Gaji:         9500000,
					FotoKTP:      "URL FOTO KTP",
					FotoSelfie:   "URL FOTO SELFIE",
					Limit:        20000000,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal(err)
			}

			sqlxDb := sqlx.NewDb(db, "sqlmock")
			c := repository.NewCustomerRepo(sqlxDb)

			m.ExpectBegin()

			query := "INSERT INTO customers (nik, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp, foto_selfie, limit_transaction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

			expectQuery := m.ExpectExec(query).WithArgs(
				tt.args.req.Nik,
				tt.args.req.FullName,
				tt.args.req.LegalName,
				tt.args.req.TempatLahir,
				tt.args.req.TanggalLahir,
				tt.args.req.Gaji,
				tt.args.req.FotoKTP,
				tt.args.req.FotoSelfie,
				tt.args.req.Limit,
			)

			if tt.wantErr != nil {
				expectQuery.WillReturnError(tt.wantErr)
			} else {
				expectQuery.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			tx, _ := sqlxDb.BeginTxx(tt.args.ctx, nil)

			gotErr := c.InsertCustomer(tt.args.ctx, tx, tt.args.req)

			assert.Equal(t, tt.wantErr, gotErr)

		})
	}
}

func Test_customerRepo_CheckLimitTransaction(t *testing.T) {
	ctx := context.Background()
	type mockExec struct {
		data int
	}
	type args struct {
		ctx        context.Context
		otrPrice   int
		customerId int64
	}
	tests := []struct {
		name     string
		args     args
		want     bool
		mockExec mockExec
		wantErr  error
	}{
		{
			name: "best_case",
			args: args{
				ctx:      ctx,
				otrPrice: 1000000,
			},
			mockExec: mockExec{
				data: 1,
			},
			want:    true,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal(err)
			}

			sqlxDb := sqlx.NewDb(db, "sqlmock")

			m.ExpectBegin()

			query := "SELECT limit_transaction > ? as check_limit FROM customers where id = ? FOR UPDATE"

			expectQuery := m.ExpectQuery(query).WithArgs(tt.args.otrPrice, tt.args.customerId)

			if tt.wantErr != nil {
				expectQuery.WillReturnError(tt.wantErr)
			} else {
				row := sqlmock.NewRows([]string{"check_limit"})
				row.AddRow(tt.mockExec.data)
				expectQuery.WillReturnRows(row)
			}

			tx, _ := sqlxDb.BeginTxx(tt.args.ctx, nil)

			c := repository.NewCustomerRepo(sqlxDb)
			got, gotErr := c.CheckLimitTransaction(tt.args.ctx, tx, tt.args.otrPrice, tt.args.customerId)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

//
// func Test_customerRepo_UpdateLimitTransaction(t *testing.T) {
// 	type fields struct {
// 		baseRepo baseRepo
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		tx         TxProvider
// 		otrPrice   int
// 		customerId int64
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &customerRepo{
// 				baseRepo: tt.fields.baseRepo,
// 			}
// 			if err := c.UpdateLimitTransaction(tt.args.ctx, tt.args.tx, tt.args.otrPrice, tt.args.customerId); (err != nil) != tt.wantErr {
// 				t.Errorf("customerRepo.UpdateLimitTransaction() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//
// func Test_customerRepo_InsertContract(t *testing.T) {
// 	type fields struct {
// 		baseRepo baseRepo
// 	}
// 	type args struct {
// 		ctx context.Context
// 		tx  TxProvider
// 		req request.CreateContactRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    int64
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &customerRepo{
// 				baseRepo: tt.fields.baseRepo,
// 			}
// 			got, err := c.InsertContract(tt.args.ctx, tt.args.tx, tt.args.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("customerRepo.InsertContract() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("customerRepo.InsertContract() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
