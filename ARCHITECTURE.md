# Архитектура HealthHub Bot API

## Обзор

Проект использует **Domain-Driven Design (DDD)** архитектуру с разделением на слои:

```
┌─────────────────────────────────────┐
│   Presentation Layer (GraphQL)     │
│   - Resolvers                      │
│   - GraphQL Schema                 │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Application Layer (Use Cases)    │
│   - CreateSymptomUseCase           │
│   - GenerateReportUseCase          │
│   - ...                             │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Domain Layer                      │
│   - Entities                        │
│   - Repositories (interfaces)       │
│   - Domain Services                 │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Infrastructure Layer              │
│   - Repository implementations      │
│   - Database (GORM)                 │
│   - File Storage                    │
│   - External Services               │
└─────────────────────────────────────┘
```

## Структура проекта

```
health-hub-bot-api/
├── cmd/
│   └── server/              # Точка входа приложения
│       └── main.go
├── internal/
│   ├── domain/              # Домены (DDD)
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   └── repository.go
│   │   ├── symptom/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── errors.go
│   │   ├── analysis/
│   │   ├── medication/
│   │   └── doctorvisit/
│   ├── application/         # Use Cases (Application Services)
│   │   ├── symptom/
│   │   │   └── create_symptom.go
│   │   └── doctorvisit/
│   │       └── generate_report.go
│   ├── infrastructure/      # Инфраструктура
│   │   ├── database/
│   │   │   └── postgres.go
│   │   ├── repository/
│   │   │   ├── user_repository.go
│   │   │   ├── symptom_repository.go
│   │   │   └── ...
│   │   └── storage/
│   │       └── file_storage.go
│   └── presentation/        # GraphQL Layer
│       └── graphql/
│           └── resolver.go
├── graphql/
│   ├── schema.graphql       # GraphQL схема
│   └── generated/           # Сгенерированный код
├── migrations/              # Миграции БД
└── config/                 # Конфигурация
```

## Домены

### 1. User Domain
**Ответственность**: Управление пользователями и их профилями

**Сущности**:
- `User` — пользователь системы

**Реpository**: `user.Repository`

### 2. Symptom Domain
**Ответственность**: Медицинский дневник (симптомы, самочувствие, показатели)

**Сущности**:
- `SymptomEntry` — запись симптома

**Repository**: `symptom.Repository`

**Use Cases**:
- `CreateSymptomUseCase` — создание записи симптома
- `GetSymptomTrendUseCase` — получение тренда самочувствия

### 3. Analysis Domain
**Ответственность**: Хранилище анализов (фото/PDF, группировка, напоминания)

**Сущности**:
- `Analysis` — медицинский анализ

**Repository**: `analysis.Repository`

**Use Cases**:
- `CreateAnalysisUseCase` — создание анализа
- `GetAnalysesByTypeUseCase` — группировка по типам

### 4. Medication Domain
**Ответственность**: Учёт лекарств (дозировки, приём, напоминания)

**Сущности**:
- `Medication` — лекарство
- `MedicationIntake` — факт приёма лекарства

**Repositories**:
- `medication.Repository` — лекарства
- `medication.IntakeRepository` — приёмы

**Use Cases**:
- `CreateMedicationUseCase` — создание лекарства
- `MarkIntakeUseCase` — отметка приёма
- `GetComplianceRateUseCase` — процент соблюдения режима

### 5. DoctorVisit Domain
**Ответственность**: Подготовка к визиту (история, динамика, вопросы)

**Сущности**:
- `DoctorVisit` — визит к врачу
- `Report` — отчёт для визита

**Repository**: `doctorvisit.Repository`

**Use Cases**:
- `CreateDoctorVisitUseCase` — создание визита
- `GenerateReportUseCase` — генерация отчёта

## Принципы DDD

### 1. Агрегаты
Каждый домен имеет свой агрегат:
- `User` — корневой агрегат для User Domain
- `SymptomEntry` — корневой агрегат для Symptom Domain
- `Analysis` — корневой агрегат для Analysis Domain
- `Medication` — корневой агрегат для Medication Domain
- `DoctorVisit` — корневой агрегат для DoctorVisit Domain

### 2. Репозитории
Репозитории определены как интерфейсы в domain слое, реализация — в infrastructure слое.

### 3. Use Cases
Use Cases находятся в application слое и координируют работу между доменами.

### 4. Инварианты
- `SymptomEntry.WellbeingScale` должен быть от 1 до 10
- `Medication` должен иметь корректное расписание
- `DoctorVisit` должен принадлежать пользователю

## GraphQL API

### Запросы (Queries)
- `me` — текущий пользователь
- `symptoms` — список симптомов с фильтрацией
- `analyses` — список анализов
- `medications` — список лекарств
- `doctorVisits` — список визитов
- `doctorVisitReport` — отчёт для визита

### Мутации (Mutations)
- `updateUserProfile` — обновление профиля
- `createSymptomEntry` — создание записи симптома
- `createAnalysis` — создание анализа
- `createMedication` — создание лекарства
- `markMedicationIntake` — отметка приёма
- `createDoctorVisit` — создание визита
- `generateDoctorVisitReport` — генерация отчёта

## База данных

### Технологии
- **PostgreSQL** — основная БД
- **GORM** — ORM для работы с БД

### Миграции
Миграции будут храниться в `migrations/` и выполняться через инструмент миграций (например, golang-migrate или GORM migrations).

## Аутентификация

### Telegram Mini App
Пользователи аутентифицируются через Telegram WebApp:
1. Telegram передаёт `initData` в заголовке запроса
2. Сервер проверяет подпись `initData`
3. Извлекает `user.id` из данных
4. Создаёт или получает пользователя по `telegram_user_id`

### Контекст
После аутентификации `user.ID` добавляется в контекст GraphQL запроса.

## Файловое хранилище

### Требования
- Хранение фото симптомов
- Хранение PDF/фото анализов

### Реализация
- **MVP**: Локальное хранилище (`./storage/`)
- **Production**: S3 или совместимое хранилище

## Напоминания

### Механизм
- Периодическая задача (cron job) проверяет предстоящие напоминания
- Отправка через Telegram Bot API
- Статус напоминания сохраняется в БД

### Типы напоминаний
1. Приём лекарств (по расписанию)
2. Сдача анализов (следующий анализ)
3. Проверка самочувствия (ежедневно)

## Обработка ошибок

### Domain Errors
Ошибки домена определяются в каждом домене (`errors.go`):
- `ErrSymptomNotFound`
- `ErrInvalidWellbeingScale`
- и т.д.

### GraphQL Errors
Domain ошибки преобразуются в GraphQL ошибки в resolvers.

## Тестирование

### Unit Tests
- Тесты для domain entities
- Тесты для use cases
- Тесты для repositories (с моками)

### Integration Tests
- Тесты GraphQL resolvers
- Тесты с реальной БД (testcontainers)

## Развёртывание

### Локальная разработка
```bash
make docker-up      # Запуск PostgreSQL
make migrate        # Применение миграций
make run            # Запуск сервера
```

### Production
- Docker контейнер с приложением
- PostgreSQL в отдельном контейнере или managed service
- Миграции через CI/CD pipeline

## Безопасность

### Защита данных
- Шифрование данных в покое (БД)
- HTTPS для всех запросов
- Валидация всех входных данных

### Telegram Security
- Проверка подписи `initData`
- Rate limiting для API запросов
- Защита от SQL injection (GORM)

## Мониторинг и логирование

### Логирование
- Structured logging (logrus или zap)
- Уровни: DEBUG, INFO, WARN, ERROR

### Метрики
- Количество запросов к GraphQL
- Время выполнения запросов
- Ошибки по типам

## Масштабирование

### Горизонтальное масштабирование
- Stateless API серверы
- Shared PostgreSQL database
- Shared file storage (S3)

### Оптимизация
- Индексы БД для частых запросов
- Кэширование (Redis) для часто запрашиваемых данных
- Пагинация для больших списков

