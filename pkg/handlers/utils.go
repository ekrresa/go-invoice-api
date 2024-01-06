package handlers

import "github.com/ekrresa/invoice-api/pkg/models"

func mergeInvoices(dest *models.UpdateInvoiceInput, src *models.Invoice) {
	if dest.Description == nil {
		dest.Description = src.Description
	}
	if dest.Status == nil {
		dest.Status = &src.Status
	}
	if dest.CustomerName == nil {
		dest.CustomerName = &src.CustomerName
	}
	if dest.CustomerEmail == nil {
		dest.CustomerEmail = src.CustomerEmail
	}
	if dest.AllowMultiplePayments == nil {
		dest.AllowMultiplePayments = &src.AllowMultiplePayments
	}
	if dest.Currency == nil {
		dest.Currency = &src.Currency
	}
	if dest.Total == nil {
		dest.Total = &src.Total
	}
	if dest.DueDate == nil {
		dest.DueDate = src.DueDate
	}
}
