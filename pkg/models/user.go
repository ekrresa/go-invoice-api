package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"name" json:"name"`
	Email     string    `gorm:"unique;uniqueIndex" json:"email"`
	Password  string    `gorm:"password" json:"password"`
	ApiKey    string    `gorm:"api_key" json:"api_key"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}
