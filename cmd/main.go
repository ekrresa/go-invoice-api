package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ekrresa/invoice-api/pkg/config"
	"github.com/ekrresa/invoice-api/pkg/helpers"
	"github.com/ekrresa/invoice-api/pkg/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var PORT = helpers.GetEnv("PORT")

	db := config.ConnectToDatabase()
	config.ApplyMigrations(db)

	var router = chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	// Rate limit by IP address and endpoint
	router.Use(httprate.Limit(
		10,
		5*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Rate limit by API key
	router.Use(httprate.Limit(
		100,
		1*time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			return r.Header.Get("X-API-Key"), nil
		}),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			helpers.ErrorResponse(w, "Too many requests", http.StatusTooManyRequests)
		})))

	router.Use(middleware.RequestSize(1048576))
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)

	routes.RegisterRoutes(router, db)

	var server = &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Printf("Server starting on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
