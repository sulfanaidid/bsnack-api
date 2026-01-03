package main

import (
	"fmt"
	"log"
	"os"

	"bsnack/internal/config"
	"bsnack/internal/handler"
	"bsnack/internal/repository"
	"bsnack/internal/router"
	"bsnack/internal/service"
	"bsnack/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbConfig := config.LoadDatabaseConfig()
	fmt.Println("dbConfig.DSN()", dbConfig.DSN())
	db, err := database.NewPostgres(dbConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	trxRepo := repository.NewTransactionRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)

	svc := &service.TransactionService{
		DB:              db,
		CustomerRepo:    customerRepo,
		ProductRepo:     productRepo,
		TransactionRepo: trxRepo,
	}

	h := &handler.TransactionHandler{Service: svc}
	r := router.Setup(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5004"
	}

	addr := ":" + port
	log.Printf("Server running on %s\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
