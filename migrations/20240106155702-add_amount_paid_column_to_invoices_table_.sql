
-- +migrate Up
ALTER TABLE invoices
ADD COLUMN amount_paid INT NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE invoices
DROP COLUMN amount_paid;