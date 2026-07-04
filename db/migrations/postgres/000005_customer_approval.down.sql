DROP INDEX IF EXISTS idx_customer_status;
DROP INDEX IF EXISTS idx_customer_status_pending;

ALTER TABLE customer
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS rejected_by,
    DROP COLUMN IF EXISTS rejected_at,
    DROP COLUMN IF EXISTS remark;