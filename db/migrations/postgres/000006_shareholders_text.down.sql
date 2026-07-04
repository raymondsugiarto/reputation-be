-- Reverting 000006 widens the column back to varchar(255). Because
-- some rows may now exceed 255 chars, the conversion truncates them
-- — this is intentional so the down migration succeeds without
-- losing the rest of the column.
ALTER TABLE customer
    ALTER COLUMN shareholders TYPE varchar(255)
        USING LEFT(shareholders, 255);