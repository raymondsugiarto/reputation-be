-- 000005_customer_approval
--
-- Adds the approval lifecycle columns to the customer table. Every new
-- customer sign-up lands in PENDING_APPROVAL and can only be moved to
-- APPROVED or REJECTED by an internal admin via /api/admin/approvals.
--
-- Backfill:
--   Existing customer rows (from before the approval workflow) are
--   marked as APPROVED so the legacy seeded `owner` user can still log
--   in. Without this backfill, every pre-existing customer would be
--   blocked at sign-in.

ALTER TABLE customer
    ADD COLUMN IF NOT EXISTS status      varchar(50)  NOT NULL DEFAULT 'PENDING_APPROVAL',
    ADD COLUMN IF NOT EXISTS approved_by varchar(255) NULL,
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP    NULL,
    ADD COLUMN IF NOT EXISTS rejected_by varchar(255) NULL,
    ADD COLUMN IF NOT EXISTS rejected_at TIMESTAMP    NULL,
    ADD COLUMN IF NOT EXISTS remark      text         NULL;

-- Backfill existing rows so the legacy `owner` user (and any customers
-- already created before this migration) are not blocked from signing
-- in. The IF NOT EXISTS-style default above already handles new rows,
-- so this UPDATE only touches pre-existing rows.
UPDATE customer
   SET status      = 'APPROVED',
       approved_by = 'system',
       approved_at = NOW()
 WHERE status = 'PENDING_APPROVAL';

-- Index for the admin approval queue. Status is the most common filter
-- (PENDING_APPROVAL) so a partial index on that value is the cheapest
-- way to back the queue query.
CREATE INDEX IF NOT EXISTS idx_customer_status_pending
    ON customer (created_at DESC)
 WHERE status = 'PENDING_APPROVAL';

CREATE INDEX IF NOT EXISTS idx_customer_status
    ON customer (status);