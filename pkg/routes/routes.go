package routes

import (
	"github.com/ekrresa/invoice-api/pkg/controllers"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(r *chi.Mux, db *gorm.DB) {
	repo := repository.NewRepository(db)
	ctrl := controllers.NewUserController(repo)

	r.Post("/users", ctrl.RegisterUser)
}
