package symptom

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Filter представляет фильтр для поиска записей симптомов
type Filter struct {
	UserID          uuid.UUID
	StartDate       *time.Time
	EndDate         *time.Time
	MinWellbeingScale *int
	MaxWellbeingScale *int
}

// Repository определяет интерфейс для работы с записями симптомов
type Repository interface {
	// Create создаёт новую запись симптома
	Create(ctx context.Context, entry *SymptomEntry) error
	
	// GetByID возвращает запись по ID
	GetByID(ctx context.Context, id uuid.UUID) (*SymptomEntry, error)
	
	// FindByFilter возвращает записи по фильтру
	FindByFilter(ctx context.Context, filter Filter, limit, offset int) ([]*SymptomEntry, int, error)
	
	// Update обновляет запись
	Update(ctx context.Context, entry *SymptomEntry) error
	
	// Delete удаляет запись
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetWellbeingTrend возвращает тренд самочувствия за период
	GetWellbeingTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]WellbeingDataPoint, error)
}

// WellbeingDataPoint представляет точку данных для графика самочувствия
type WellbeingDataPoint struct {
	Date  time.Time
	Value int
}

