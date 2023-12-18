package models

type CreateInvoicePayload struct {
	Description           string `json:"description,omitempty"`
	Status                string `json:"status,omitempty"`
	AllowMultiplePayments bool   `json:"allow_multiple_payments,omitempty"`
	CustomerName          string `json:"customer_name" validate:"required"`
	CustomerEmail         string `json:"customer_email,omitempty"`
	Currency              string `json:"currency" validate:"required"`
	Total                 uint   `json:"total" validate:"required"`
	DueDate               string `json:"due_date" validate:"required"`
	Items                 []struct {
		Name     string `json:"name" validate:"required"`
		Quantity uint   `json:"quantity" validate:"required"`
		Price    uint   `json:"price" validate:"required"`
	} `json:"items,omitempty"`
}

type InvoiceResponse struct {
	ID                    string        `json:"id"`
	Description           uint          `json:"description"`
	Status                InvoiceStatus `json:"status"`
	AllowMultiplePayments bool          `json:"allow_multiple_payments"`
	Items                 []InvoiceItem `json:"items"`
	Total                 uint          `json:"total"`
}
