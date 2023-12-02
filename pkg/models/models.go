package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:50" json:"id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Email     string    `gorm:"not null;size:255;unique;uniqueIndex" json:"email"`
	Password  string    `gorm:"not null;size:255" json:"password"`
	ApiKey    string    `gorm:"not null;size:255" json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Invoices  []Invoice `json:"invoices"`
}

type Invoice struct {
	ID            string        `gorm:"primaryKey;size:50" json:"id"`
	UserID        string        `gorm:"user_id;size:50" json:"user_id"`
	Description   string        `json:"description"`
	Status        InvoiceStatus `gorm:"default:draft" json:"status"`
	CustomerName  string        `gorm:"customer_name;not null;size:255" json:"customer_name"`
	CustomerEmail string        `gorm:"customer_email;size:255" json:"customer_email"`
	Underpay      bool          `gorm:"default:false" json:"underpay"`
	Currency      string        `gorm:"default:NGN;size:5" json:"currency"`
	Total         uint64        `gorm:"not null" json:"total"`
	DueDate       time.Time     `gorm:"not null" json:"due_date"`
	CreatedAt     time.Time     `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"updated_at" json:"updated_at"`
}

type InvoiceStatus string

const (
	Open    InvoiceStatus = "open"
	Draft   InvoiceStatus = "draft"
	Paid    InvoiceStatus = "paid"
	Void    InvoiceStatus = "void"
	Deleted InvoiceStatus = "deleted"
)

type InvoiceItem struct {
	ID        uint      `gorm:"autoIncrement;primary_key" json:"id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Quantity  uint      `gorm:"default:1" json:"quantity"`
	Price     uint      `gorm:"not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	InvoiceID string    `gorm:"not null;size:50" json:"invoice_id"`
	Invoice   Invoice   `json:"invoice"`
}
