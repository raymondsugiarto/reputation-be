
CREATE TABLE IF NOT EXISTS organization (
    id varchar(255) PRIMARY KEY,
    code varchar(100) NULL,
    name varchar(255) NULL,
    origin varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "user" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_type varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "user_credential" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_id varchar(100) NULL,
    username varchar(510) NULL,
    password varchar(1020) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS customer (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_id varchar(255),
    nama_lengkap varchar(255),
    nomor_ktp varchar(255),
    nomor_npwp varchar(255),
    tanggal_lahir TIMESTAMP,
    kota_lahir varchar(255),
    status_pernikahan varchar(50),
    pendidikan_terakhir varchar(50),
    lama_tinggal varchar(255),
    alamat_jalan varchar(255),
    kecamatan varchar(255),
    kota_kabupaten varchar(255),
    provinsi varchar(255),
    kode_pos varchar(20),
    sama_dengan_domisili boolean,
    alamat_ktp_jalan varchar(255),
    kecamatan_ktp varchar(255),
    kota_kabupaten_ktp varchar(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS account (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    customer_id varchar(255),
    account_name varchar(255),
    account_type varchar(50),
    balance numeric(20, 4),
    currency varchar(10),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS user_relationship (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255),
    user_id_follower varchar(255),
    relationship varchar(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);