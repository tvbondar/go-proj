// Главная точка запуска (инициализация, сервер, consumer)
// Главная точка запуска (инициализация, сервер, consumer)
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"sync"
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

	// Конструктор кеша теперь может вернуть ошибку
	cacheRepo, err := repositories.NewCacheOrderRepository()
	if err != nil {
		log.Fatal("Failed to create cache repository:", err)
	}

	// Загружаем кеш (не блокирующую операцию можно вынести в фон)
	if err := cacheRepo.LoadFromDB(context.Background(), dbRepo); err != nil {
		log.Printf("Failed to load cache from DB (non-fatal): %v", err)
	}

	processUseCase := usecases.NewProcessOrderUseCase(dbRepo, cacheRepo)
	getUseCase := usecases.NewGetOrderUseCase(cacheRepo, dbRepo)

	// Создание контекста для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Ждём завершения горутин
	var wg sync.WaitGroup

	// Запуск Kafka Consumer в горутине с контекстом; передаём cfg, чтобы consumer не читал конфиг сам
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := handlers.StartKafkaConsumer(ctx, cfg, processUseCase); err != nil {
			log.Printf("Kafka consumer stopped with error: %v", err)
		} else {
			log.Println("Kafka consumer stopped")
		}
	}()

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

	// дождёмся завершения фоновых ворутин (например, Kafka consumer)
	wg.Wait()

	// все defer (включая db.Close()) выполнятся, выходим нормально
	log.Println("Graceful shutdown completed, exiting...")
}
