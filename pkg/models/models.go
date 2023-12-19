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
	UserID                string        `gorm"index" json:"user_id" db:"user_id"`
	Description           string        `json:"description"`
	Status                InvoiceStatus `gorm:"default:draft" json:"status"`
	CustomerName          string        `gorm:"customer_name;not null;size:255" json:"customer_name"`
	CustomerEmail         string        `gorm:"customer_email;size:255" json:"customer_email"`
	AllowMultiplePayments bool          `gorm:"default:false" json:"allow_multiple_payments"`
	Currency              string        `gorm:"default:NGN;size:5" json:"currency"`
	Total                 uint          `gorm:"not null" json:"total"`
	DueDate               time.Time     `gorm:"not null" json:"due_date"`
	CreatedAt             time.Time     `gorm:"created_at" json:"created_at"`
	UpdatedAt             time.Time     `gorm:"updated_at" json:"updated_at"`
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
