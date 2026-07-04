-- 000004_customer_company_fields
--
-- Adds the company-specific columns required by the perusahaan (COMPANY)
-- sign-up flow. These columns are populated by customer.Service.SignUp only
-- when customer_type = 'COMPANY'; for INDIVIDUAL rows they remain NULL.
--
-- Naming is in English snake_case (matching the rest of the schema) instead
-- of the Indonesian field names used by the FE.

ALTER TABLE customer
    ADD COLUMN IF NOT EXISTS company_name        varchar(255)     NULL,
    ADD COLUMN IF NOT EXISTS establishment_date  TIMESTAMP        NULL,
    ADD COLUMN IF NOT EXISTS business_sector     varchar(255)     NULL,
    ADD COLUMN IF NOT EXISTS kbli_code           varchar(255)     NULL,
    ADD COLUMN IF NOT EXISTS company_tax_id      varchar(255)     NULL,
    ADD COLUMN IF NOT EXISTS shareholders        varchar(255)     NULL;

-- The application serialises the shareholder slice to JSON in Go
-- and stores it as a string. We use varchar(255) here for the initial
-- shape so the migration stays additive and reversible; migration
-- 000006 widens this column to `text` to fit realistic shareholder
-- lists. See db/migrations/postgres/000006_shareholders_text.up.sql.