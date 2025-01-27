package main

import (
	"context"
	"finance-api/internal/handler"
	"finance-api/internal/repository"
	"finance-api/internal/service"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	// Получаем конфигурацию из ENV
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", 
        os.Getenv("DB_USER"), 
        os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_HOST"), 
        os.Getenv("DB_PORT"), 
        os.Getenv("DB_NAME"), 
        os.Getenv("DB_SSLMODE"))

    // Настройка пула соединений
    dbPool, err := pgxpool.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatal("Failed to connect to the database: ", err)
    }
    defer dbPool.Close()

    // Инициализация репозиториев
    userRepo := repository.NewUserRepository(dbPool)
    transactionRepo := repository.NewTransactionRepository(dbPool)

    // Инициализация сервисов
    userService := service.NewUserService(userRepo)
    transactionService := service.NewTransactionService(transactionRepo)

    // Инициализация обработчиков
    userHandler := &handler.UserHandler{Service: userService}
    transactionHandler := handler.NewTransactionHandler(transactionService)

    // Настройка маршрутов Gin
    r := gin.Default()

    // Маршруты для пользователей
    r.POST("/users", userHandler.CreateUser)
    r.POST("/users/balance", userHandler.AddBalance)
    r.POST("/users/transfer", userHandler.TransferBalance)
	r.GET("/users/:id", userHandler.GetUserByID)

    // Маршрут для получения последних транзакций
    r.GET("/transactions", transactionHandler.GetLastTransactions)

	// Получаем порт из ENV, если он не задан - используем 8080
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to run server: ", err)
    }
}
