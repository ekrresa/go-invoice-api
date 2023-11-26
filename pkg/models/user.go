package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique;uniqueIndex" json:"email"`
	Password  string    `json:"password"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
