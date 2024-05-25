CREATE TABLE IF NOT EXISTS customers (
	id BIGINT NOT NULL AUTO_INCREMENT,
	nik varchar(16) NOT NULL,
	full_name varchar(30) NOT NULL,
	legal_name varchar(30) NOT NULL,
	tempat_lahir varchar(30) NOT NULL,
	tanggal_lahir varchar(30) NOT NULL,
	gaji INT NOT NULL,
	foto_ktp VARCHAR(255) NOT NULL,
	foto_selfie VARCHAR(255) NOT NULL,
	limit_transaction INT NOT NULL,

	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS contracts (
	id BIGINT NOT NULL AUTO_INCREMENT,
	customer_id BIGINT NOT NULL,
	no_kontrak VARCHAR(255) NOT NULL,
	otr_price INT NOT NULL,
	tenor INT NOT NULL,
	bunga INT NOT NULL,
	nama_aset VARCHAR(255) NOT NULL,
	fee INT NOT NULL,

	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS installments (
	id BIGINT NOT NULL AUTO_INCREMENT,
	contract_id BIGINT NOT NULL,
	installment_amount INT NOT NULL,
	paid_amount INT DEFAULT 0,
	due_date TIMESTAMP NOT NULL,

	PRIMARY KEY (id)
);






