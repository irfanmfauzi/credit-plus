CREATE TABLE customers (
  id bigint NOT NULL AUTO_INCREMENT,
  nik varchar(16) NOT NULL,
  full_name varchar(30) NOT NULL,
  legal_name varchar(30) NOT NULL,
  tempat_lahir varchar(30) NOT NULL,
  tanggal_lahir varchar(30) NOT NULL,
  gaji int NOT NULL,
  foto_ktp varchar(255) NOT NULL,
  foto_selfie varchar(255) NOT NULL,
  limit_transaction int NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE contracts (
  id bigint NOT NULL AUTO_INCREMENT,
  customer_id bigint NOT NULL,
  no_kontrak varchar(255) NOT NULL,
  otr_price int NOT NULL,
  tenor int NOT NULL,
  bunga int NOT NULL,
  nama_aset varchar(255) NOT NULL,
  fee int NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE installments (
  id bigint NOT NULL AUTO_INCREMENT,
  contract_id bigint NOT NULL,
  installment_amount int NOT NULL,
  paid_amount int DEFAULT '0',
  due_date timestamp NOT NULL,
  PRIMARY KEY (`id`)
)
