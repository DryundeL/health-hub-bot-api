package medication

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository определяет интерфейс для работы с лекарствами
type Repository interface {
	// Create создаёт новое лекарство
	Create(ctx context.Context, medication *Medication) error
	
	// GetByID возвращает лекарство по ID
	GetByID(ctx context.Context, id uuid.UUID) (*Medication, error)
	
	// FindByUserID возвращает все лекарства пользователя
	FindByUserID(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]*Medication, error)
	
	// Update обновляет лекарство
	Update(ctx context.Context, medication *Medication) error
	
	// Delete удаляет лекарство
	Delete(ctx context.Context, id uuid.UUID) error
}

// IntakeRepository определяет интерфейс для работы с приёмами лекарств
type IntakeRepository interface {
	// Create создаёт запись о приёме лекарства
	Create(ctx context.Context, intake *MedicationIntake) error
	
	// GetByID возвращает запись по ID
	GetByID(ctx context.Context, id uuid.UUID) (*MedicationIntake, error)
	
	// FindByMedicationAndDate возвращает приёмы за конкретную дату
	FindByMedicationAndDate(ctx context.Context, medicationID uuid.UUID, date time.Time) ([]*MedicationIntake, error)
	
	// Update обновляет запись о приёме
	Update(ctx context.Context, intake *MedicationIntake) error
	
	// GetUpcomingIntakes возвращает предстоящие приёмы
	GetUpcomingIntakes(ctx context.Context, userID uuid.UUID, fromTime time.Time, limit int) ([]*MedicationIntake, error)
	
	// GetComplianceRate возвращает процент соблюдения режима приёма
	GetComplianceRate(ctx context.Context, medicationID uuid.UUID, startDate, endDate time.Time) (float64, error)
}

