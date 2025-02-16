package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/avito-shop-service/internal/config"
	"github.com/avito-shop-service/internal/handlers"
	"github.com/avito-shop-service/internal/repository"
	"github.com/avito-shop-service/internal/router"
	"github.com/avito-shop-service/internal/services"
)

func main() {
	testMode := os.Getenv("TEST_MODE") == "true"

	cfg, err := config.LoadConfig(testMode)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}
	dbConfig.MaxConns = 50

	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	merchRepo := repository.NewMerchRepository(db)

	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	inventoryService := services.NewInventoryService(inventoryRepo)
	merchService := services.NewMerchService(merchRepo)

	userHandler := handlers.NewUserHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(userService, transactionService)
	buyHandler := handlers.NewBuyHandler(userService, merchService, inventoryService, transactionService)
	infoHandler := handlers.NewInformationHandler(userService, merchService, inventoryService, transactionService)

	r := router.NewRouter(transactionHandler, userHandler, buyHandler, infoHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server is running port: %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
