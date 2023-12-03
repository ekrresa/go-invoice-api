package middleware

import (
	"net/http"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/ekrresa/invoice-api/pkg/utils"
)

type Middleware struct {
	db *repository.Repository
}

func NewMiddleware(db *repository.Repository) *Middleware {
	return &Middleware{
		db: db,
	}
}

func (m *Middleware) AuthenticateApiKey(fn func(http.ResponseWriter, *http.Request, *models.UserWithoutPassword)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			utils.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		apiKeyHash := utils.HashApiKey(apiKey)

		user, userErr := m.db.GetUserByApiKey(apiKeyHash)
		if userErr != nil {
			utils.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fn(w, r, user)
	}
}
