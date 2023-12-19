package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func ConnectToDatabase() *sqlx.DB {
	connString, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Fatal("DATABASE_URL not set")
	}

	var db, err = sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	fmt.Println("Connected to DB!")

	return db
}

func ApplyMigrations(db *sqlx.DB) {
	var migrations = &migrate.FileMigrationSource{Dir: "migrations"}
	var n, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	if err != nil {
		log.Fatal("Error applying migrations:", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}
