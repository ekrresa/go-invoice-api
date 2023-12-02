package config

import (
	"errors"
	"os"

	"github.com/ekrresa/invoice-api/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	connUrl, exists := os.LookupEnv("CONN_STRING")
	if !exists {
		return nil, errors.New("CONN_STRING not set")
	}

	db, err := gorm.Open(postgres.Open(connUrl), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{}, &models.Invoice{}, &models.InvoiceItem{})

	return db, nil
}
