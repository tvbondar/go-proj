// Главная точка запуска (инициализация, сервер, consumer)
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/tvbondar/go-server/internal/config"
	"github.com/tvbondar/go-server/internal/handlers"
	"github.com/tvbondar/go-server/internal/repositories"
	"github.com/tvbondar/go-server/internal/usecases"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	var db *sql.DB
	// Попытка подключения к БД с миграциями
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", cfg.DBDSN)
		if err != nil {
			log.Printf("Failed to open database connection (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		// Применение миграций перед ping
		if err := goose.Up(db, "migrations"); err != nil {
			log.Printf("Failed to apply migrations (attempt %d): %v", i+1, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping database (attempt %d): %v", i+1, err)
		db.Close()
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database and applied migrations")

	dbRepo := repositories.NewPostgresOrderRepository(db)
	cacheRepo := repositories.NewCacheOrderRepository()

	if err := cacheRepo.LoadFromDB(dbRepo); err != nil {
		log.Fatal("Failed to load cache from DB:", err)
	}

	processUseCase := usecases.NewProcessOrderUseCase(dbRepo, cacheRepo)
	getUseCase := usecases.NewGetOrderUseCase(cacheRepo, dbRepo)

	// Создание контекста для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск Kafka Consumer в горутине с контекстом
	go handlers.StartKafkaConsumer(ctx, processUseCase)

	httpHandler := handlers.NewHTTPHandler(getUseCase)
	mux := http.NewServeMux()
	mux.HandleFunc("/order/", httpHandler.GetOrder)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mux,
	}

	// Запуск HTTP-сервера в горутине
	go func() {
		log.Printf("HTTP server starting on %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("Shutdown signal received, starting graceful shutdown...")

	// Остановка сервера с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	// Закрытие БД
	log.Println("Graceful shutdown completed, exiting...")
	os.Exit(0)
}
