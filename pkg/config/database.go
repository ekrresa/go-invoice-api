package config

import (
	"fmt"
	"log"

	"github.com/ekrresa/invoice-api/pkg/helpers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func ConnectToDatabase() *sqlx.DB {
	var DB_URL = helpers.GetEnv("DATABASE_URL")

	var db, err = sqlx.Connect("postgres", DB_URL)
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
