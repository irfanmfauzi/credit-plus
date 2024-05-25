package request

import "time"

type CreateCustomerRequest struct {
	Nik          string `json:"nik"`
	FullName     string `json:"full_name"`
	LegalName    string `json:"legal_name"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Gaji         int    `json:"gaji"`
	FotoKTP      string `json:"foto_ktp"`
	FotoSelfie   string `json:"foto_selfie"`
	Limit        int    `json:"limit"`
}

type CreateContactRequest struct {
	CustomerID int64  `json:"customer_id"`
	NoKontrak  string `json:"no_kontrak"`
	OtrPrice   int    `json:"otr_price"`
	Tenor      int    `json:"tenor"`
	Bunga      int    `json:"bunga"`
	NamaAset   string `json:"nama_aset"`
	Fee        int    `json:"fee"`
}

type CreateCustomerInstallmentRequest struct {
	ContractId        int64     `db:"contract_id"`
	InstallmentAmount int       `db:"installment_amount"`
	PaidAmount        int       `db:"paid_amount"`
	DueDate           time.Time `db:"due_date"`
}
