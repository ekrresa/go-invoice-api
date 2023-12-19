package repository

import (
	"fmt"
	"log"
	"strings"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/oklog/ulid/v2"
)

func (r *Repository) CreateInvoice(userID string, newInvoice *models.CreateInvoiceInput) error {
	var invoiceID = strings.ToLower(ulid.Make().String())

	var tx, txErr = r.db.Beginx()
	if txErr != nil {
		return txErr
	}

	var _, invoiceInsertErr = tx.Exec(`INSERT INTO invoices (id, user_id, description, status, customer_name, customer_email, allow_multiple_payments, currency, total, due_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		invoiceID, userID, newInvoice.Description, newInvoice.Status, newInvoice.CustomerName, newInvoice.CustomerEmail, newInvoice.AllowMultiplePayments, newInvoice.Currency, newInvoice.Total, newInvoice.DueDate)

	if invoiceInsertErr != nil {
		log.Println(invoiceInsertErr.Error(), "Line 23")
		tx.Rollback()
		return invoiceInsertErr
	}

	if len(newInvoice.Items) > 0 {
		for index := range newInvoice.Items {
			newInvoice.Items[index].InvoiceID = invoiceID
		}

		fmt.Println(newInvoice.Items)

		var _, invoiceItemsInsertErr = tx.NamedExec(`INSERT INTO invoice_items (invoice_id, name, quantity, unit_price)
		 VALUES (:invoice_id, :name, :quantity, :unit_price)`, newInvoice.Items)

		if invoiceItemsInsertErr != nil {
			log.Println(invoiceItemsInsertErr.Error(), "Line 37")
			tx.Rollback()
			return invoiceItemsInsertErr
		}
	}

	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return txCommitErr
	}

	return nil
}

// TODO: include invoice items as an option
func (r *Repository) ListInvoicesOfUser(userID string) ([]models.Invoice, error) {
	var invoices = []models.Invoice{}

	var err = r.db.Select(&invoices, `SELECT * FROM invoices WHERE user_id = $1`, userID)

	return invoices, err
}

// TODO: Include the invoice items
func (r *Repository) GetInvoice(invoiceID string, userID string) (models.Invoice, error) {
	var invoice models.Invoice

	var err = r.db.Get(&invoice, `SELECT * FROM invoices WHERE id = $1 AND user_id = $2`,
		invoiceID, userID)

	return invoice, err
}
