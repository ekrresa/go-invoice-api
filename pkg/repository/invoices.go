package repository

import (
	"github.com/ekrresa/invoice-api/pkg/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateInvoice(invoice *models.Invoice, invoiceItems *[]models.InvoiceItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}

		if invoiceItems != nil && len(*invoiceItems) > 0 {
			if err := tx.Model(&models.InvoiceItem{}).Create(invoiceItems).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) ListInvoicesOfUser(userID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	error := r.db.Table("invoices").Where("user_id = ?", userID).Find(&invoices).Error

	return invoices, error
}

// func (r *Repository) GetInvoice(userID string) ([]models.Invoice, error) {
// 	var invoices []models.Invoice
// 	error := r.db.Table("invoices").Where("user_id = ?", userID).Find(&invoices).Error

// 	return invoices, error
// }
