CREATE TABLE IF NOT EXISTS funder (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255),
    name varchar(255),
    phone_number varchar(255),
    funder_id_parent varchar(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS contract (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    funder_id varchar(255),
    contract_number int NULL,
    contract_code varchar(255) NULL,
    disbursement_at TIMESTAMP,
    amount numeric(20, 4),
    total_paid_amount numeric(20, 4) NULL,
    duration int NULL,
    due_date TIMESTAMP NULL,
    destination_account varchar(510) NULL,
    return_percentage numeric(5, 2) NULL,
    return_amount numeric(20, 4) NULL,
    attachment_url varchar(255) NULL,
    notes varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS contract_payment (
    id varchar(255) PRIMARY KEY,
    contract_id varchar(255),
    payment_at TIMESTAMP,
    payment_amount numeric(20, 4),
    attachment_url varchar(255) NULL,
    notes varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
