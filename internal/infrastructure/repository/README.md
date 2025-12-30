# Репозитории

Этот пакет содержит реализации репозиториев для работы с PostgreSQL через GORM.

## Структура

- `user_repository.go` - репозиторий пользователей
- `symptom_repository.go` - репозиторий записей симптомов
- `analysis_repository.go` - репозиторий анализов
- `medication_repository.go` - репозиторий лекарств
- `medication_intake_repository.go` - репозиторий приёмов лекарств
- `doctor_visit_repository.go` - репозиторий визитов к врачу

## Использование

```go
import (
    "github.com/health-hub-bot-api/internal/infrastructure/database"
    "github.com/health-hub-bot-api/internal/infrastructure/repository"
)

// Подключение к БД
db, err := database.NewPostgresFromEnv()
if err != nil {
    log.Fatal(err)
}
defer database.Close(db)

// Создание репозиториев
userRepo := repository.NewUserRepository(db)
symptomRepo := repository.NewSymptomRepository(db)
// ...
```

## Особенности реализации

- Все репозитории используют GORM для работы с PostgreSQL
- Модели БД отделены от доменных сущностей
- Поддержка контекста для отмены операций
- Пагинация для больших списков
- Мягкое удаление (soft delete) для пользователей

