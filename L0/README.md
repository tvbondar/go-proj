# Wildberries L0 – Демонстрационный сервис с Kafka, PostgreSQL, кешем


## Архитектура системы

```mermaid
flowchart TD
    A[Web Client<br/>(HTML + JS)] 
    -->|HTTP GET /order/:id| B(Go Service<br/>http://localhost:8081)

    B --> C{Kafka Message<br/>Received?}
    C -->|Yes| D[Kafka Consumer<br/>(goroutine)]
    D -->|Parse JSON| E[Validate & Save]
    E --> F[(PostgreSQL)]
    E --> G[In-memory Cache<br/>(sync.Map / LRU)]

    B --> H{Cache Hit?}
    H -->|Yes| I[Return from Cache<br/>⚡ Instant]
    H -->|No| J[Load from PostgreSQL]
    J --> G
    J --> I

    K[Kafka Broker<br/>(orders topic)] --> D

    style B fill:#4CAF50,stroke:#333,color:white
    style G fill:#2196F3,stroke:#333,color:white
    style F fill:#FF9800,stroke:#333,color:white
    style I fill:#8BC34A,stroke:#333,color:white


---

### Как запустить (30 секунд до рабочего сервиса)

```bash
git clone https://github.com/tvbondar/go-proj/L0.git
cd L0
# Запуск всего стека
task compose-up          
# Отправка тестового заказа
task publish-order
# Открываем
http://localhost:8081

```

Вводим b563feb7b2b84b6test, чтобы увидеть сохраненный заказ

## Полезные команды (Taskfile)

|Команда             |   Описание                                                  |
|--------------------|-------------------------------------------------------------|
|task compose-up     | Поднять весь стек                                            |
|task publish-order  | Отправить  заказ в Kafka                                     | 
|task restart-app    | Перезапустить приложение, показать восстановление кэша из БД |
|task db-clean       | Очистить все заказы                                          |
|task test           | Запустить все тесты                                          |
|task logs           | Следить за логами приложения                                 |


## Технологии и инструменты
**Backend**
1. Go 1.24
2. net/http 
3. segmentio/kafka-go – Kafka-клиент
4. PostgreSQL 15 + goose миграции
5. hashicorp/golang-lru/v2 – LRU кэш 
6. validator.v10 – валидация структур
7. zap – структурированное логирование
8. gomock -  мокки 

**Инфраструктура**
Docker + Docker Compose
Kafka + Zookeeper (Confluent)
PostgreSQL
Taskfile 

**Frontend**
1. HTML 

## Production-конфигурация

```bash
cp configs/.env.example configs/.env.prod
nano configs/.env.prod
```

**GitHub**: @tvbondar
**ТГ**: @Sinopah
 



