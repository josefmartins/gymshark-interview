package main

import (
	"context"
	"errors"
	"fmt"
	"gymshark-interview/database/migrations"
	"gymshark-interview/internal/server"
	"gymshark-interview/internal/service"
	"gymshark-interview/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	defaultHTTPServerPort  = 8080
	shutdownGracefulPeriod = 2 * time.Second
)

func main() {
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

	// initialise dependencies
	repo := storage.New(db)
	productService := service.NewProductService(repo)
	packageService := service.NewPackageService(repo)

	port := defaultHTTPServerPort
	portFromEnv := os.Getenv("SERVER_PORT")
	if len(portFromEnv) > 0 {
		port, err = strconv.Atoi(portFromEnv)
		if err != nil {
			log.Fatal("SERVER_PORT value is not valid: ", err)
		}
	}

	server := server.New(port, productService, packageService)

	// start server
	go server.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// shutdown server
	log.Println("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownGracefulPeriod)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("server shutdown gracefully")
}
