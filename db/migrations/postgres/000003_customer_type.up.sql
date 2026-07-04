ALTER TABLE customer
    ADD COLUMN IF NOT EXISTS customer_type varchar(50) NULL;

-- Nullable because the column is new and pre-existing rows have no type.
-- Going forward, application code (customer.Service.SignUp) requires the value.