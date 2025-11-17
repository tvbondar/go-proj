# go-server

Л0/
├── cmd/                    # Точка входа приложения
│   └── main.go            # Главная точка запуска (инициализация, сервер, consumer)
├── internal/               # Логика приложения (приватный код)
│   ├── entities/          # Модели данных (структуры заказов)
│   │   └── order.go       # Структура Order, Delivery, Payment, Items
|   ├── logger/
|   │   └── logger.go
│   ├── handlers/          # Обработчики HTTP и Kafka
│   │   ├── http_handler.go # HTTP-эндпоинты (Gin)
│   │   └── kafka_handler.go # Логика обработки сообщений из Kafka
│   ├── repositories/      # Абстракции для доступа к данным
|   │   ├── mocks/
|   |   │   └── mocks_order_repository.go
│   │   ├── cache_repository.go # Интерфейс и реализация In-memory Cache
│   │   ├── postgres_repository.go # Работа с PostgreSQL
│   │   └── repository.go  # Интерфейс Repository (общий для кэша и БД)
│   ├── usecases/          # Бизнес-логика (Service)
│   │   ├── get_order.go   # Логика получения заказа (из кэша/БД)
|   |   ├── get_order_test.go
|   |   ├── process_order_test.go
│   │   └── process_order.go # Логика обработки нового заказа (сохранение)
│   └── config/            # Конфигурации
│       └── config.go      # Загрузка настроек (из env или YAML)
├── web/                    # Статические файлы фронтенда
│   └── index.html         # Простой веб-интерфейс
├── migrations/            # Миграции БД
│   └── 001_create_tables.sql # Одна миграция для таблиц
├── Dockerfile             # Конфигурация Docker
├── docker-compose.yml     # Оркестрация (Go, Kafka, PostgreSQL)
├── go.mod                 # Зависимости проекта
├── go.sum                 # Хэши зависимостей
└── model.json             # Модель данных заказов (для справки)


