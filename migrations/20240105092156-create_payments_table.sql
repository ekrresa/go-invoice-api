
-- +migrate Up
CREATE TABLE IF NOT EXISTS payments (
  id VARCHAR(100) PRIMARY KEY,
  invoice_id VARCHAR(100) NOT NULL,
	amount INT NOT NULL,
  reference VARCHAR(50) NOT NULL,
	customer_email VARCHAR(255) DEFAULT NULL,
  currency CHAR(3) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_invoice_payment FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
);

CREATE INDEX idx_invoice_payment ON payments (invoice_id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE PROCEDURE auto_set_timestamp();

-- +migrate Down
DROP TRIGGER IF EXISTS set_timestamp ON payments;
DROP INDEX IF EXISTS idx_invoice_payment;
DROP TABLE IF EXISTS payments;