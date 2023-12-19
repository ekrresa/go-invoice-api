package middleware

import (
	"net/http"

	"github.com/ekrresa/invoice-api/pkg/helpers"
	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
)

type Middleware struct {
	db *repository.Repository
}

func NewMiddleware(db *repository.Repository) *Middleware {
	return &Middleware{
		db: db,
	}
}

func (m *Middleware) AuthenticateApiKey(fn func(http.ResponseWriter, *http.Request, *models.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			helpers.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var currentUser models.User

		var cachedUser, found = helpers.GetUserFromCache(apiKey)

		if !found {
			apiKeyHash := helpers.HashApiKey(apiKey)
			user, userErr := m.db.GetUserByApiKey(apiKeyHash)

			if userErr != nil {
				helpers.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			helpers.CacheUser(apiKey, user)
			currentUser = user
		} else {
			currentUser = cachedUser
		}

		fn(w, r, &currentUser)
	}
}
