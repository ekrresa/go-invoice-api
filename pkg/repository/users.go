package repository

import (
	"github.com/ekrresa/invoice-api/pkg/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
}

func (r *Repository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	return err
}
