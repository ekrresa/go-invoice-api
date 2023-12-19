package routes

import (
	"github.com/ekrresa/invoice-api/pkg/handlers"
	"github.com/ekrresa/invoice-api/pkg/middleware"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *chi.Mux, db *sqlx.DB) {
	repo := repository.NewRepository(db)
	middleware := middleware.NewMiddleware(repo)
	userHandler := handlers.NewUserHandler(*repo)
	invoiceHandler := handlers.NewInvoiceHandler(*repo)

	r.Post("/users/auth", userHandler.RegisterUser)
	r.Post("/users/get_apikey", userHandler.RegenerateApiKey)

	r.Post("/invoices", middleware.AuthenticateApiKey(invoiceHandler.CreateInvoice))
	r.Get("/invoices", middleware.AuthenticateApiKey(invoiceHandler.ListInvoicesOfUser))
	r.Get("/invoices/{invoiceID}", middleware.AuthenticateApiKey(invoiceHandler.GetInvoice))
}
