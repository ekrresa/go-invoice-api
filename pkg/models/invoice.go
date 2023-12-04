package models

type CreateInvoicePayload struct {
	Description   string `json:"description,omitempty"`
	Status        string `json:"status,omitempty"`
	Underpay      bool   `json:"underpay,omitempty"`
	CustomerName  string `json:"customer_name" validate:"required"`
	CustomerEmail string `json:"customer_email,omitempty"`
	Currency      string `json:"currency" validate:"required"`
	Total         uint   `json:"total" validate:"required"`
	DueDate       string `json:"due_date" validate:"required"`
	Items         []struct {
		Name     string `json:"name" validate:"required"`
		Quantity uint   `json:"quantity" validate:"required"`
		Price    uint   `json:"price" validate:"required"`
	} `json:"items,omitempty"`
}
