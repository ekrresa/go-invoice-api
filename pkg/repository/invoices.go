package repository

import (
	"github.com/ekrresa/invoice-api/pkg/models"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	CreateInvoice() error
}

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
