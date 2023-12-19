package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ekrresa/invoice-api/pkg/config"
	"github.com/ekrresa/invoice-api/pkg/routes"
	"github.com/ekrresa/invoice-api/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	db := config.ConnectToDatabase()
	config.ApplyMigrations(db)

	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// Rate limit by IP address and endpoint
	r.Use(httprate.Limit(
		10,
		5*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Rate limit by API key
	r.Use(httprate.Limit(
		100,
		1*time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			return r.Header.Get("X-API-Key"), nil
		}),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			utils.ErrorResponse(w, "Too many requests", http.StatusTooManyRequests)
		})))

	r.Use(middleware.RequestSize(1048576))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	routes.RegisterRoutes(r, db)

	err := http.ListenAndServe(":8000", r)

	if err != nil {
		log.Fatal(err)
	}
}
