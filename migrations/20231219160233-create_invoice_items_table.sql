
-- +migrate Up
CREATE TABLE IF NOT EXISTS invoice_items (
  id SERIAL PRIMARY KEY,
  invoice_id VARCHAR(100) NOT NULL,
  name VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  unit_price INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id)
);

CREATE INDEX idx_invoice_items_id ON invoice_items (invoice_id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON invoice_items
FOR EACH ROW
EXECUTE PROCEDURE auto_set_timestamp();

-- +migrate Down
DROP TRIGGER IF EXISTS set_timestamp ON invoice_items;
DROP INDEX IF EXISTS idx_invoice_items_id;
DROP TABLE IF EXISTS invoice_items;
