package repository

import (
	"github.com/ekrresa/invoice-api/pkg/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

func (r *Repository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	return err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Select("ID", "Name", "Email", "Password", "ApiKey", "CreatedAt", "UpdatedAt").
		Where("email = ?", email).
		First(&user).Error
	return user, err
}

func (r *Repository) GetUserByApiKey(apiKey string) (*models.UserWithoutPassword, error) {
	user := &models.UserWithoutPassword{}
	err := r.db.
		Select("ID", "Name", "Email", "ApiKey", "CreatedAt", "UpdatedAt").
		Where("api_key = ?", apiKey).
		First(&user).Error

	return user, err
}

func (r *Repository) UpdateUser(user *models.User) error {
	err := r.db.Save(user).Error
	return err
}
