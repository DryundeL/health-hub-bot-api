package doctorvisit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository определяет интерфейс для работы с визитами к врачу
type Repository interface {
	// Create создаёт новый визит
	Create(ctx context.Context, visit *DoctorVisit) error
	
	// GetByID возвращает визит по ID
	GetByID(ctx context.Context, id uuid.UUID) (*DoctorVisit, error)
	
	// FindByUserID возвращает визиты пользователя
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*DoctorVisit, int, error)
	
	// Update обновляет визит
	Update(ctx context.Context, visit *DoctorVisit) error
	
	// Delete удаляет визит
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetUpcomingVisits возвращает предстоящие визиты
	GetUpcomingVisits(ctx context.Context, userID uuid.UUID, beforeDate time.Time) ([]*DoctorVisit, error)
}

