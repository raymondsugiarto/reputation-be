-- 000006_shareholders_text
--
-- Widen the `shareholders` column from varchar(255) to text.
--
-- Background:
--   000004_customer_company_fields created `shareholders` as
--   varchar(255). The Go model declares the column with a jsonb GORM
--   tag, but PostgreSQL stored it as a plain string regardless —
--   the column type and the GORM tag are independent. varchar(255) is
--   a 255-char ceiling, which is far too small for a realistic
--   shareholder list (a single Indonesian company with full KBLI /
--   NPWP data and two shareholders already serialises past 255).
--
--   The migration runner applies `*.up.sql` in numeric order. This
--   migration runs after 000005_customer_approval, so any seed or
--   backfill that referenced the column under its old shape is
--   already in place.
--
-- We use `text` rather than `jsonb` on purpose: the application
-- serialises the shareholder slice to JSON via Go and stores it as
-- a string. The Go layer parses it on read. Promoting the column to
-- `jsonb` would force every write path to issue a real JSON parse,
-- which is unnecessary for a small bounded struct like
-- PemegangSahamDto. Keeping it as `text` is the smallest change
-- that removes the overflow without altering the runtime contract.

ALTER TABLE customer
    ALTER COLUMN shareholders TYPE text
        USING shareholders::text;