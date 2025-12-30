.PHONY: help generate run test docker-up docker-down migrate

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

generate: ## Сгенерировать GraphQL код
	go run github.com/99designs/gqlgen generate

run: ## Запустить сервер
	go run cmd/server/main.go

test: ## Запустить тесты
	go test ./...

docker-up: ## Запустить PostgreSQL в Docker
	docker-compose up -d

docker-down: ## Остановить PostgreSQL в Docker
	docker-compose down

migrate: ## Запустить миграции (TODO: добавить реализацию)
	@echo "TODO: Добавить команду для миграций"

install: ## Установить зависимости
	go mod download
	go mod tidy

