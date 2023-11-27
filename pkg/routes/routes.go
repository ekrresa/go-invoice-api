package routes

import (
	"github.com/ekrresa/invoice-api/pkg/handlers"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(r *chi.Mux, db *gorm.DB) {
	repo := repository.NewRepository(db)
	userHandler := handlers.NewUserController(repo)

	r.Post("/users", userHandler.RegisterUser)
}
