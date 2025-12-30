# HealthHub Bot API

API для Telegram mini-app HealthHub — персонального координатора здоровья.

## Архитектура

Проект использует:
- **DDD (Domain-Driven Design)** — разделение на домены
- **GraphQL** — API через gqlgen
- **PostgreSQL** — основная БД
- **Redis** — кэширование и очереди
- **GORM** — ORM для работы с БД

## Структура проекта

```
health-hub-bot-api/
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── domain/          # Домены (DDD)
│   │   ├── user/
│   │   ├── symptom/
│   │   ├── analysis/
│   │   ├── medication/
│   │   └── doctorvisit/
│   ├── infrastructure/  # Инфраструктура (БД, внешние сервисы)
│   ├── application/     # Use cases / Application services
│   └── presentation/    # GraphQL resolvers, handlers
├── graphql/
│   ├── schema.graphql   # GraphQL схема
│   └── generated/       # Сгенерированный код gqlgen
├── migrations/          # Миграции БД
└── config/             # Конфигурация
```

## Запуск

### Требования
- Go 1.21+ (для локальной разработки)
- Docker и Docker Compose (для запуска через Docker)
- PostgreSQL 14+ (если запускаете локально)
- Redis 7+ (если запускаете локально)

### Запуск через Docker (рекомендуется)

#### Полный запуск (приложение + БД + Redis)
```bash
# Сборка и запуск всех сервисов
docker-compose up -d --build

# Просмотр логов
docker-compose logs -f app

# Остановка всех сервисов
docker-compose down

# Остановка с удалением volumes (удалит данные БД и Redis)
docker-compose down -v
```

#### Запуск только инфраструктуры (для локальной разработки)
```bash
# Запуск PostgreSQL и Redis
docker-compose -f docker-compose.dev.yml up -d

# Остановка
docker-compose -f docker-compose.dev.yml down
```

После запуска инфраструктуры приложение можно запускать локально:
```bash
go run cmd/server/main.go
```

### Локальная разработка

#### Установка зависимостей
```bash
go mod download
```

#### Настройка окружения
Создайте файл `.env` на основе `env.example`:

**Вариант 1: Использование отдельных переменных (рекомендуется)**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=healthhub
DB_PASSWORD=healthhub
DB_NAME=healthhub
DB_SSLMODE=disable
DB_TIMEZONE=UTC
PORT=8080
TELEGRAM_BOT_TOKEN=your_bot_token
```

**Вариант 2: Использование DATABASE_URL**
```env
DATABASE_URL=postgres://healthhub:healthhub@localhost:5432/healthhub?sslmode=disable
PORT=8080
TELEGRAM_BOT_TOKEN=your_bot_token
```

#### Генерация GraphQL кода
```bash
go run github.com/99designs/gqlgen generate
```

#### Запуск миграций
```bash
# TODO: добавить команду для миграций
```

#### Запуск сервера
```bash
go run cmd/server/main.go
```

### Docker команды

```bash
# Пересборка образа приложения
docker-compose build app

# Перезапуск приложения
docker-compose restart app

# Просмотр логов конкретного сервиса
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f app

# Выполнение команд в контейнере
docker-compose exec app sh
docker-compose exec postgres psql -U healthhub -d healthhub
docker-compose exec redis redis-cli
```

## Разработка

### Добавление нового домена
1. Создайте директорию в `internal/domain/{domain_name}/`
2. Определите entities, repositories, services
3. Добавьте GraphQL типы в `graphql/schema.graphql`
4. Создайте resolvers в `internal/presentation/graphql/`

### Тестирование
```bash
go test ./...
```

## Документация

См. [PRODUCT.md](./PRODUCT.md) для продуктового описания MVP.

