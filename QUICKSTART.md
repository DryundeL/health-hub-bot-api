# Быстрый старт HealthHub Bot API

## Предварительные требования

- Go 1.21 или выше
- PostgreSQL 14+
- Docker (опционально, для локальной БД)

## Шаг 1: Клонирование и установка зависимостей

```bash
# Установка зависимостей
make install
# или
go mod download
go mod tidy
```

## Шаг 2: Настройка базы данных

### Вариант A: Использование Docker (рекомендуется для разработки)

```bash
# Запуск PostgreSQL
make docker-up

# Применение миграций
psql -h localhost -U healthhub -d healthhub -f migrations/001_initial_schema.sql
# Пароль: healthhub
```

### Вариант B: Использование существующей PostgreSQL

1. Создайте базу данных:
```sql
CREATE DATABASE healthhub;
```

2. Примените миграции:
```bash
psql -U your_user -d healthhub -f migrations/001_initial_schema.sql
```

## Шаг 3: Настройка окружения

Создайте файл `.env` (скопируйте из `.env.example`):

```env
DATABASE_URL=postgres://healthhub:healthhub@localhost:5432/healthhub?sslmode=disable
PORT=8080
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

## Шаг 4: Генерация GraphQL кода

```bash
make generate
# или
go run github.com/99designs/gqlgen generate
```

## Шаг 5: Запуск сервера

```bash
make run
# или
go run cmd/server/main.go
```

Сервер будет доступен по адресу: `http://localhost:8080`

GraphQL Playground: `http://localhost:8080/`

## Проверка работы

Откройте GraphQL Playground и выполните тестовый запрос:

```graphql
query {
  __typename
}
```

## Следующие шаги

1. **Реализация репозиториев**: Создайте реализации репозиториев в `internal/infrastructure/repository/`
2. **Реализация use cases**: Завершите реализацию use cases в `internal/application/`
3. **Реализация resolvers**: Завершите реализацию GraphQL resolvers
4. **Аутентификация**: Добавьте middleware для проверки Telegram WebApp данных
5. **Файловое хранилище**: Реализуйте сохранение фото/PDF файлов

## Структура проекта

```
health-hub-bot-api/
├── cmd/server/          # Точка входа
├── internal/
│   ├── domain/          # Домены (DDD)
│   ├── application/     # Use cases
│   ├── infrastructure/  # Репозитории, БД, хранилище
│   └── presentation/    # GraphQL resolvers
├── graphql/             # GraphQL схема и сгенерированный код
├── migrations/          # SQL миграции
└── config/             # Конфигурация
```

## Полезные команды

```bash
make help          # Показать все доступные команды
make generate      # Сгенерировать GraphQL код
make run           # Запустить сервер
make test          # Запустить тесты
make docker-up     # Запустить PostgreSQL
make docker-down   # Остановить PostgreSQL
```

## Документация

- [PRODUCT.md](./PRODUCT.md) — продуктовый документ MVP
- [ARCHITECTURE.md](./ARCHITECTURE.md) — архитектура проекта
- [README.md](./README.md) — общая информация о проекте

## Разработка

### Добавление нового домена

1. Создайте директорию в `internal/domain/{domain_name}/`
2. Определите entity, repository interface, errors
3. Добавьте GraphQL типы в `graphql/schema.graphql`
4. Создайте use cases в `internal/application/{domain_name}/`
5. Реализуйте repository в `internal/infrastructure/repository/`
6. Реализуйте resolver в `internal/presentation/graphql/`
7. Запустите `make generate` для генерации GraphQL кода

### Тестирование

```bash
# Unit тесты
go test ./internal/domain/...

# Integration тесты
go test ./internal/infrastructure/...

# Все тесты
make test
```

## Troubleshooting

### Ошибка подключения к БД

Проверьте:
- PostgreSQL запущен
- Правильность `DATABASE_URL` в `.env`
- Доступность порта 5432

### Ошибки генерации GraphQL

Убедитесь, что:
- Все типы в `schema.graphql` корректны
- Запущен `make generate`
- Нет синтаксических ошибок в схеме

### Проблемы с зависимостями

```bash
go mod tidy
go mod download
```

