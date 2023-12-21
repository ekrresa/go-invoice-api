package models

import (
	"encoding/json"
	"time"
)

type CreateInvoiceItemInput struct {
	Name      string `json:"name" validate:"required"`
	Quantity  uint   `json:"quantity" validate:"required"`
	UnitPrice uint   `json:"unit_price" validate:"required" db:"unit_price"`
	InvoiceID string `json:"invoice_id,omitempty" db:"invoice_id"`
}

type CreateInvoiceInput struct {
	Description           *string                  `json:"description,omitempty"`
	Status                InvoiceStatus            `json:"status,omitempty"`
	AllowMultiplePayments bool                     `json:"allow_multiple_payments,omitempty"`
	CustomerName          string                   `json:"customer_name" validate:"required"`
	CustomerEmail         *string                  `json:"customer_email,omitempty"`
	Currency              string                   `json:"currency" validate:"required"`
	Total                 uint                     `json:"total" validate:"required"`
	DueDate               *time.Time               `json:"due_date,omitempty"`
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

type GetInvoiceResponse struct {
	ID                    string                   `json:"id"`
	Description           *string                  `json:"description"`
	Status                InvoiceStatus            `json:"status"`
	AllowMultiplePayments bool                     `json:"allow_multiple_payments" db:"allow_multiple_payments"`
	CustomerName          string                   `json:"customer_name"`
	CustomerEmail         *string                  `json:"customer_email"`
	Currency              string                   `json:"currency"`
	DueDate               *time.Time               `json:"due_date" db:"due_date"`
	CreatedAt             time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at" db:"updated_at"`
	Total                 uint                     `json:"total"`
	Items                 []GetInvoiceItemResponse `json:"items"`
}

type GetInvoiceItemResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Quantity  uint      `json:"quantity"`
	UnitPrice uint      `json:"unit_price" db:"unit_price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateInvoiceInput struct {
	Description           *string        `json:"description,omitempty"`
	Status                *InvoiceStatus `json:"status,omitempty"`
	AllowMultiplePayments *bool          `json:"allow_multiple_payments,omitempty"`
	// TODO: Add validation for customer email
	CustomerName  *string    `json:"customer_name,omitempty"`
	CustomerEmail *string    `json:"customer_email,omitempty"`
	Currency      *string    `json:"currency,omitempty"`
	Total         *uint      `json:"total,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
}
