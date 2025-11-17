// Главная точка запуска (инициализация, сервер, consumer)
// Главная точка запуска (инициализация, сервер, consumer)
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"strings"
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

func ensureAddr(port string) string {
	if port == "" {
		return ":8081"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var db *sql.DB
	var lastErr error
	connected := false

	// Попытка подключения к БД с миграциями
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", cfg.DBDSN)
		if err != nil {
			lastErr = err
			log.Printf("Failed to open database connection (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Убедимся, что база отвечает
		if err = db.Ping(); err != nil {
			lastErr = err
			log.Printf("Failed to ping database (attempt %d): %v", i+1, err)
			if cerr := db.Close(); cerr != nil {
				log.Printf("failed to close db: %v", cerr)
			}
			time.Sleep(2 * time.Second)
			continue
		}

		// Применение миграций после успешного ping
		if err := goose.Up(db, "migrations"); err != nil {
			lastErr = err
			log.Printf("Failed to apply migrations (attempt %d): %v", i+1, err)
			if cerr := db.Close(); cerr != nil {
				log.Printf("failed to close db: %v", cerr)
			}
			time.Sleep(2 * time.Second)
			continue
		}

		// Всё успешно
		connected = true
		break
	}

	if !connected {
		log.Fatalf("Failed to connect to database after retries: %v", lastErr)
	}

	// Закрываем DB только при завершении приложения
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	log.Println("Successfully connected to database and applied migrations")

	dbRepo := repositories.NewPostgresOrderRepository(db)

	// Конструктор кеша теперь может вернуть ошибку
	cacheRepo, err := repositories.NewCacheOrderRepository()
	if err != nil {
		log.Fatalf("Failed to create cache repository: %v", err)
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
		Addr:    ensureAddr(cfg.HTTPPort),
		Handler: mux,
		// можно добавить таймауты при необходимости
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Запуск HTTP-сервера в горутине
	go func() {
		log.Printf("HTTP server starting on %s", server.Addr)
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
