package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository определяет интерфейс для работы с пользователями
type Repository interface {
	// Create создаёт нового пользователя
	Create(ctx context.Context, user *User) error

	// GetByID возвращает пользователя по ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByTelegramUserID возвращает пользователя по Telegram User ID
	GetByTelegramUserID(ctx context.Context, telegramUserID int64) (*User, error)

	// Update обновляет пользователя
	Update(ctx context.Context, user *User) error

	// Delete удаляет пользователя (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error
}
