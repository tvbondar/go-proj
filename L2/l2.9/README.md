## Задание L2.9
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