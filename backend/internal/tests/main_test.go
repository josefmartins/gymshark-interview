package tests

import (
	"context"
	"fmt"
	"gymshark-interview/database/migrations"
	"gymshark-interview/internal/server"
	"gymshark-interview/internal/service"
	"gymshark-interview/internal/storage"
	"log"
	"strconv"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

var hostname string

func TestMain(m *testing.M) {
	// start in-memory sqlite with empty db
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// run migrations + seeds
	mg := migrations.GetMigrationSource()
	migrate.SetTable("migrations")
	_, err = migrate.Exec(db.DB, "sqlite3", mg, migrate.Up)
	if err != nil {
		log.Fatal(fmt.Errorf("migrations failed: %w", err))
	}
	// init deps + server
	repo := storage.New(db)
	packageService := service.NewPackageService(repo)
	productService := service.NewProductService(repo)

	port := 3000
	server := server.New(port, productService, packageService)

	hostname = "http://localhost:" + strconv.Itoa(port)

	go server.Start()

	_ = m.Run()

	_ = server.Shutdown(context.Background())
}
