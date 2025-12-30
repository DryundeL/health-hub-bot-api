package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/medication"
	"gorm.io/gorm"
)

// medicationIntakeModel представляет модель приёма лекарства в БД
type medicationIntakeModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	MedicationID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	ScheduledTime time.Time  `gorm:"not null;index"`
	TakenAt       *time.Time
	IsTaken       bool       `gorm:"not null;default:false;index"`
	Notes         *string    `gorm:"type:text"`
	CreatedAt     time.Time  `gorm:"not null"`
}

// TableName возвращает имя таблицы
func (medicationIntakeModel) TableName() string {
	return "medication_intakes"
}

// toDomain преобразует модель БД в доменную сущность
func (m *medicationIntakeModel) toDomain() *medication.MedicationIntake {
	return &medication.MedicationIntake{
		ID:            m.ID,
		MedicationID:  m.MedicationID,
		ScheduledTime: m.ScheduledTime,
		TakenAt:       m.TakenAt,
		IsTaken:       m.IsTaken,
		Notes:         m.Notes,
		CreatedAt:     m.CreatedAt,
	}
}

// fromDomain преобразует доменную сущность в модель БД
func (m *medicationIntakeModel) fromDomain(intake *medication.MedicationIntake) {
	m.ID = intake.ID
	m.MedicationID = intake.MedicationID
	m.ScheduledTime = intake.ScheduledTime
	m.TakenAt = intake.TakenAt
	m.IsTaken = intake.IsTaken
	m.Notes = intake.Notes
	m.CreatedAt = intake.CreatedAt
}

// IntakeRepository реализует medication.IntakeRepository для PostgreSQL
type IntakeRepository struct {
	db *gorm.DB
}

// NewIntakeRepository создаёт новый репозиторий приёмов лекарств
func NewIntakeRepository(db *gorm.DB) medication.IntakeRepository {
	return &IntakeRepository{db: db}
}

// Create создаёт запись о приёме лекарства
func (r *IntakeRepository) Create(ctx context.Context, intake *medication.MedicationIntake) error {
	model := &medicationIntakeModel{}
	model.fromDomain(intake)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*intake = *model.toDomain()
	return nil
}

// GetByID возвращает запись по ID
func (r *IntakeRepository) GetByID(ctx context.Context, id uuid.UUID) (*medication.MedicationIntake, error) {
	var model medicationIntakeModel
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return model.toDomain(), nil
}

// FindByMedicationAndDate возвращает приёмы за конкретную дату
func (r *IntakeRepository) FindByMedicationAndDate(ctx context.Context, medicationID uuid.UUID, date time.Time) ([]*medication.MedicationIntake, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var models []medicationIntakeModel
	if err := r.db.WithContext(ctx).
		Where("medication_id = ? AND scheduled_time >= ? AND scheduled_time < ?", medicationID, startOfDay, endOfDay).
		Order("scheduled_time ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	intakes := make([]*medication.MedicationIntake, len(models))
	for i := range models {
		intakes[i] = models[i].toDomain()
	}

	return intakes, nil
}

// Update обновляет запись о приёме
func (r *IntakeRepository) Update(ctx context.Context, intake *medication.MedicationIntake) error {
	model := &medicationIntakeModel{}
	model.fromDomain(intake)

	return r.db.WithContext(ctx).
		Model(&medicationIntakeModel{}).
		Where("id = ?", intake.ID).
		Updates(model).Error
}

// GetUpcomingIntakes возвращает предстоящие приёмы
func (r *IntakeRepository) GetUpcomingIntakes(ctx context.Context, userID uuid.UUID, fromTime time.Time, limit int) ([]*medication.MedicationIntake, error) {
	var models []medicationIntakeModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN medications ON medications.id = medication_intakes.medication_id").
		Where("medications.user_id = ? AND medication_intakes.scheduled_time >= ? AND medication_intakes.is_taken = ?", userID, fromTime, false).
		Order("medication_intakes.scheduled_time ASC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	intakes := make([]*medication.MedicationIntake, len(models))
	for i := range models {
		intakes[i] = models[i].toDomain()
	}

	return intakes, nil
}

// GetComplianceRate возвращает процент соблюдения режима приёма
func (r *IntakeRepository) GetComplianceRate(ctx context.Context, medicationID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	var result struct {
		Total int64
		Taken int64
	}

	err := r.db.WithContext(ctx).
		Model(&medicationIntakeModel{}).
		Select("COUNT(*) as total, SUM(CASE WHEN is_taken THEN 1 ELSE 0 END) as taken").
		Where("medication_id = ? AND scheduled_time >= ? AND scheduled_time <= ?", medicationID, startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	if result.Total == 0 {
		return 0, nil
	}

	return float64(result.Taken) / float64(result.Total) * 100, nil
}

