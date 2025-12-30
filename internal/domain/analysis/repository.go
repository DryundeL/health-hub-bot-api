package analysis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Filter представляет фильтр для поиска анализов
type Filter struct {
	UserID    uuid.UUID
	Type      *Type
	StartDate *time.Time
	EndDate   *time.Time
}

// Repository определяет интерфейс для работы с анализами
type Repository interface {
	// Create создаёт новый анализ
	Create(ctx context.Context, analysis *Analysis) error
	
	// GetByID возвращает анализ по ID
	GetByID(ctx context.Context, id uuid.UUID) (*Analysis, error)
	
	// FindByFilter возвращает анализы по фильтру
	FindByFilter(ctx context.Context, filter Filter, limit, offset int) ([]*Analysis, int, error)
	
	// Update обновляет анализ
	Update(ctx context.Context, analysis *Analysis) error
	
	// Delete удаляет анализ
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetByType группирует анализы по типу
	GetByType(ctx context.Context, userID uuid.UUID) (map[Type][]*Analysis, error)
	
	// GetUpcomingReminders возвращает анализы с предстоящими напоминаниями
	GetUpcomingReminders(ctx context.Context, userID uuid.UUID, beforeDate time.Time) ([]*Analysis, error)
}

