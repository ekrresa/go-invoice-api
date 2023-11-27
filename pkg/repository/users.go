package repository

import (
	"github.com/ekrresa/invoice-api/pkg/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

func (r *Repository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	return err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}
