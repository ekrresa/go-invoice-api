package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	ApiKey    string    `json:"api_key" db:"api_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type InvoiceStatus string

const (
	Open  InvoiceStatus = "open"
	Draft InvoiceStatus = "draft"
	Paid  InvoiceStatus = "paid"
	Void  InvoiceStatus = "void"
)

type Invoice struct {
	ID                    string        `json:"id"`
	UserID                string        `json:"user_id" db:"user_id"`
	Description           string        `json:"description"`
	Status                InvoiceStatus `json:"status"`
	CustomerName          string        `json:"customer_name" db:"customer_name"`
	CustomerEmail         string        `json:"customer_email" db:"customer_email"`
	AllowMultiplePayments bool          `json:"allow_multiple_payments" db:"allow_multiple_payments"`
	Currency              string        `json:"currency"`
	Total                 uint          `json:"total"`
	DueDate               time.Time     `json:"due_date" db:"due_date"`
	CreatedAt             time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at" db:"updated_at"`
}

type InvoiceItem struct {
	ID        uint      `gorm:"autoIncrement;primary_key" json:"id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Quantity  uint      `gorm:"default:1" json:"quantity"`
	Price     uint      `gorm:"not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	InvoiceID string    `gorm:"not null;size:50" json:"invoice_id"`
	Invoice   Invoice   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"invoice"`
}
