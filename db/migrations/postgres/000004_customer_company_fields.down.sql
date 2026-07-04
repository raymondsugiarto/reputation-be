ALTER TABLE customer
    DROP COLUMN IF EXISTS company_name,
    DROP COLUMN IF EXISTS establishment_date,
    DROP COLUMN IF EXISTS business_sector,
    DROP COLUMN IF EXISTS kbli_code,
    DROP COLUMN IF EXISTS company_tax_id,
    DROP COLUMN IF EXISTS shareholders;