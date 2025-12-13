# Задание L2.10
##  Задание
Текст задания можно найти в (./docs/TASK.md)

##  Решение
Используется внешняя сортировка. Реализовано:
-  Чтение из файла
-  Обязательные флаги
-  Дополнительные флаги
-  Требования к реализации (линтеры, тесты и т.п.)

### Сборка и запуск
```bash
task build
./build/sort --help
```
Проверить работу программы можно на заранее подготовленных данных:
- [`./tests`](./tests) 

Пример использования:
```bash
./build/sort -k 4 tests/tests.txt
```

### Качество кода
```bash
task install-deps # install golangci-lint to ./bin
task lint # Run linter
task vet # Run go vet
task test # Run tests
```

> [!NOTE]
> Полный список команд доступен при помощи команды `task`

> [!TIP]
> `Task` можно установить через команду:
> ```bash
> go install github.com/go-task/task/v3/cmd/task@latest
> ```


















# Инструкция по запуску решения и тестов
Список доступных команд (task --list):

* build:               Build the application
* install-tools:       Install developer tools (golangci-lint) into $(go env GOPATH)/bin
* lint:                Run golangci-lint (requires golangci-lint installed)
* list:                Print available tasks from this Taskfile
* run:                 Run the application
* test:                Run all tests
* test-unpacker:       Run only unpacker package tests (single test file / function)
* vet:                 Run go vet

Для проведения тестов используйте команды:

1. task lint - проверка через golangci-lint
2. task vet - проверка через vet
3. task test-unpacker - запуск теста в пакете unpacker
4. task test - запуск всех тестов.

Для установки зависимостей использовать команду:
 task install-tools