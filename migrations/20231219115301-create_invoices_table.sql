
-- +migrate Up
CREATE TYPE INVOICE_STATUS AS ENUM ('open', 'draft', 'paid', 'void');

CREATE TABLE IF NOT EXISTS invoices (
	id VARCHAR(100) PRIMARY KEY,
  user_id VARCHAR(100) NOT NULL,
	description TEXT DEFAULT NULL,
	status INVOICE_STATUS NOT NULL DEFAULT 'open',
	customer_name VARCHAR(255) NOT NULL,
	customer_email VARCHAR(255) DEFAULT NULL,
  allow_multiple_payments BOOLEAN NOT NULL DEFAULT FALSE,
  currency CHAR(3) NOT NULL,
	total INT NOT NULL DEFAULT 0,
  due_date TIMESTAMP DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_invoices_user_id ON invoices (user_id);
CREATE INDEX idx_invoices_status ON invoices (status);

-- +migrate Down
DROP TABLE IF EXISTS invoices;
DROP TYPE IF EXISTS INVOICE_STATUS;
DROP INDEX IF EXISTS idx_invoices_user_id;
DROP INDEX IF EXISTS idx_invoices_status;