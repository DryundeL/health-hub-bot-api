# Пакет конфигурации

Пакет `config` предоставляет централизованную загрузку и управление конфигурацией приложения.

## Использование

```go
import "github.com/health-hub-bot-api/internal/config"

// Загрузка конфигурации из переменных окружения
cfg, err := config.Load()
if err != nil {
    log.Fatal("failed to load config:", err)
}

// Использование конфигурации
db, err := database.NewPostgres(cfg.Database)
serverPort := cfg.Server.Port
```

## Структура конфигурации

### DatabaseConfig
- `Host` - хост базы данных (DB_HOST, по умолчанию localhost)
- `Port` - порт базы данных (DB_PORT, по умолчанию 5432)
- `User` - пользователь базы данных (DB_USER, обязательно если не задан DATABASE_URL)
- `Password` - пароль базы данных (DB_PASSWORD, обязательно если не задан DATABASE_URL)
- `Name` - имя базы данных (DB_NAME, обязательно если не задан DATABASE_URL)
- `SSLMode` - режим SSL подключения (DB_SSLMODE, по умолчанию disable)
- `Timezone` - часовой пояс (DB_TIMEZONE, по умолчанию UTC)
- `URL` - полная строка подключения к PostgreSQL (DATABASE_URL, альтернатива отдельным параметрам)
- `MaxOpenConns` - максимальное количество открытых соединений (DATABASE_MAX_OPEN_CONNS, по умолчанию 25)
- `MaxIdleConns` - максимальное количество неактивных соединений (DATABASE_MAX_IDLE_CONNS, по умолчанию 5)
- `LogLevel` - уровень логирования GORM (DATABASE_LOG_LEVEL: silent/error/warn/info, по умолчанию info)

### ServerConfig
- `Port` - порт сервера (PORT, по умолчанию 8080)
- `Host` - хост сервера (HOST, по умолчанию пустая строка)

### TelegramConfig
- `BotToken` - токен Telegram бота (TELEGRAM_BOT_TOKEN)

### StorageConfig
- `Type` - тип хранилища: "local" или "s3" (STORAGE_TYPE, по умолчанию "local")
- `Path` - путь для локального хранилища (STORAGE_PATH, по умолчанию "./storage")
- `S3AccessKeyID` - AWS Access Key ID (AWS_ACCESS_KEY_ID)
- `S3SecretAccessKey` - AWS Secret Access Key (AWS_SECRET_ACCESS_KEY)
- `S3Region` - AWS регион (AWS_REGION)
- `S3Bucket` - имя S3 bucket (S3_BUCKET)

## Переменные окружения

Все параметры конфигурации загружаются из переменных окружения.

### База данных

Подключение к базе данных можно настроить двумя способами:

**Способ 1: Использование DATABASE_URL (для обратной совместимости)**
```env
DATABASE_URL=postgres://user:password@localhost:5432/healthhub?sslmode=disable
```

**Способ 2: Использование отдельных переменных (рекомендуется)**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=healthhub
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

**Обязательные параметры:**
- Либо `DATABASE_URL`, либо комбинация `DB_USER`, `DB_PASSWORD`, `DB_NAME`

**Дополнительные параметры БД:**
- `DATABASE_MAX_OPEN_CONNS` - максимальное количество открытых соединений (по умолчанию 25)
- `DATABASE_MAX_IDLE_CONNS` - максимальное количество неактивных соединений (по умолчанию 5)
- `DATABASE_LOG_LEVEL` - уровень логирования GORM: silent/error/warn/info (по умолчанию info)

Остальные параметры имеют значения по умолчанию.

