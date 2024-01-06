package repository

import (
	"strings"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/jaevor/go-nanoid"
	"github.com/oklog/ulid/v2"
)

func (r *Repository) PayAnInvoice(invoice models.Invoice, input models.PayInvoiceInput, setInvoiceAsPaid bool) (*models.Payment, error) {
	var paymentID = strings.ToLower(ulid.Make().String())
	var canonicID, nanoidErr = nanoid.Standard(21)
	if nanoidErr != nil {
		return nil, nanoidErr
	}

	var payment models.Payment
	var paymentReference = strings.ToUpper(canonicID())

	var tx, txErr = r.db.Beginx()
	if txErr != nil {
		return nil, txErr
	}

	var _, paymentInsertErr = tx.Exec(`INSERT INTO payments (id, invoice_id, amount, reference, currency, customer_email) VALUES ($1, $2, $3, $4, $5, $6)`, paymentID, invoice.ID, input.Amount, paymentReference, input.Currency, input.CustomerEmail)

	if paymentInsertErr != nil {
		tx.Rollback()
		return nil, paymentInsertErr
	}

	if setInvoiceAsPaid {
		var _, invoiceInsertErr = tx.Exec(`UPDATE invoices 
		SET amount_paid = amount_paid + $1, status = 'paid' WHERE id = $2`,
			input.Amount, invoice.ID)

		if invoiceInsertErr != nil {
			tx.Rollback()
			return nil, invoiceInsertErr
		}
	} else {
		var _, invoiceInsertErr = tx.Exec(`UPDATE invoices SET amount_paid = amount_paid + $1 WHERE id = $2`, input.Amount, invoice.ID)

		if invoiceInsertErr != nil {
			tx.Rollback()
			return nil, invoiceInsertErr
		}
	}

	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return nil, txCommitErr
	}

	var getPaymentErr = r.db.Get(&payment, `SELECT * FROM payments WHERE id = $1`, paymentID)

	if getPaymentErr != nil {
		return nil, getPaymentErr
	}

	return &payment, nil
}

func (r *Repository) ListPaymentsOfAnInvoice(invoiceID string) ([]models.Payment, error) {
	var payments = []models.Payment{}

	var err = r.db.Select(&payments, `SELECT * FROM payments WHERE invoice_id = $1`, invoiceID)

	return payments, err
}
