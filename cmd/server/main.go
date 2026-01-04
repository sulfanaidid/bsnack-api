package main

import (
	"context"
	"log"
	"os"

	"bsnack/internal/config"
	"bsnack/internal/handler"
	"bsnack/internal/repository"
	"bsnack/internal/router"
	"bsnack/internal/service"
	"bsnack/pkg/database"
	pkgredis "bsnack/pkg/redis"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbConfig := config.LoadDatabaseConfig()
	db, err := database.NewPostgres(dbConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	redisClient := pkgredis.NewRedis()
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Println("redis not connected, running without cache")
	}

	transactionRepo := repository.NewTransactionRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	reportRepo := repository.NewReportRepository(db)

	transactionService := service.NewTransactionService(
		db,
		transactionRepo,
		customerRepo,
		productRepo,
	)

	reportService := service.NewReportService(
		reportRepo,
		redisClient,
	)

	transactionHandler := handler.NewTransactionHandler(transactionService)
	reportHandler := handler.NewReportHandler(reportService)

	r := router.Setup(
		transactionHandler,
		reportHandler,
	)

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
