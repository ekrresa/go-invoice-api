package repository

import (
	"strings"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/utils"
	"github.com/oklog/ulid/v2"
)

func (r *Repository) CreateUser(input models.CreateUserInput) (string, error) {
	var id = strings.ToLower(ulid.Make().String())
	var apiKey = strings.ToLower(ulid.Make().String())
	var apiKeyHash = utils.HashApiKey(apiKey)

	var _, err = r.db.Exec(`INSERT INTO users (id, name, email, password, api_key) VALUES ($1, $2, $3, $4, $5)`,
		id, input.Name, input.Email, input.Password, apiKeyHash)

	return apiKey, err
}

func (r *Repository) GetUserByEmail(email string) (models.User, error) {
	var user = models.User{}

	var err = r.db.Get(&user, `SELECT id, name, email, password, api_key 
	FROM users WHERE email = $1 LIMIT 1`, email)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByID(id string) (models.User, error) {
	var user = models.User{}

	var err = r.db.Get(&user, `SELECT id, name, email, api_key, password, created_at, updated_at 
	FROM users WHERE id = $1 LIMIT 1`, id)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByApiKey(apiKey string) (models.User, error) {
	var user = models.User{}

	var err = r.db.Get(&user, `SELECT id, name, email, api_key, password, created_at, updated_at 
	FROM users WHERE api_key = $1 LIMIT 1`, apiKey)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) UpdateUserAPIKey(userId string, apiKey string) (bool, error) {
	var result, err = r.db.Exec(`UPDATE users SET api_key = $1 WHERE id = $2`,
		apiKey, userId)

	var rows, rowsErr = result.RowsAffected()
	if rowsErr != nil {
		return false, rowsErr
	}

	return int(rows) > 0, err
}
