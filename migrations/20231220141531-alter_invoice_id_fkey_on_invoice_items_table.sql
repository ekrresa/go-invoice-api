
-- +migrate Up
ALTER TABLE invoice_items DROP CONSTRAINT fk_invoice;
ALTER TABLE invoice_items ADD CONSTRAINT fk_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) 
ON DELETE CASCADE;

-- +migrate Down
