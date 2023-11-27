package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ekrresa/invoice-api/pkg/config"
	"github.com/ekrresa/invoice-api/pkg/routes"

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
	db, dbErr := config.ConnectDatabase()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.RequestSize(1048576))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	routes.RegisterRoutes(r, db)

	err := http.ListenAndServe(":3000", r)

	if err != nil {
		log.Fatal(err)
	}
}
