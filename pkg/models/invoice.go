package models

import "encoding/json"

type CreateInvoiceItemInput struct {
	Name      string `json:"name" validate:"required"`
	Quantity  uint   `json:"quantity" validate:"required"`
	UnitPrice uint   `json:"unit_price" validate:"required" db:"unit_price"`
	InvoiceID string `json:"invoice_id,omitempty" db:"invoice_id"`
}

type CreateInvoiceInput struct {
	Description           string                   `json:"description,omitempty"`
	Status                InvoiceStatus            `json:"status,omitempty"`
	AllowMultiplePayments bool                     `json:"allow_multiple_payments,omitempty"`
	CustomerName          string                   `json:"customer_name" validate:"required"`
	CustomerEmail         string                   `json:"customer_email,omitempty"`
	Currency              string                   `json:"currency" validate:"required"`
	Total                 uint                     `json:"total" validate:"required"`
	DueDate               string                   `json:"due_date,omitempty"`
	Items                 []CreateInvoiceItemInput `json:"items,omitempty"`
}

func (i *CreateInvoiceInput) UnmarshalJSON(data []byte) error {
	type Alias CreateInvoiceInput

	var temp = Alias{
		AllowMultiplePayments: false,
		Status:                Open,
		Items:                 []CreateInvoiceItemInput{},
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*i = CreateInvoiceInput(temp)

	return nil
}

type ListInvoicesResponse struct {
	ID                    string        `json:"id"`
	Description           uint          `json:"description"`
	Status                InvoiceStatus `json:"status"`
	AllowMultiplePayments bool          `json:"allow_multiple_payments"`
	Items                 []InvoiceItem `json:"items"`
	Total                 uint          `json:"total"`
}
