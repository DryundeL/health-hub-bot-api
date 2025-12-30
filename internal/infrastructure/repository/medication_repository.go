package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/medication"
	"gorm.io/gorm"
)

// scheduleDetailsJSON представляет JSON для schedule_details
type scheduleDetailsJSON struct {
	Times []string `json:"times"`
	Days  []int    `json:"days"`
}

// Value реализует driver.Valuer для GORM
func (s scheduleDetailsJSON) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan реализует sql.Scanner для GORM
func (s *scheduleDetailsJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// medicationModel представляет модель лекарства в БД
type medicationModel struct {
	ID              uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID           `gorm:"type:uuid;not null;index"`
	Name            string              `gorm:"type:varchar(255);not null"`
	Dosage          string              `gorm:"type:varchar(100);not null"`
	ScheduleType    string              `gorm:"type:varchar(20);not null;check:schedule_type IN ('daily','weekly','as_needed')"`
	ScheduleDetails scheduleDetailsJSON `gorm:"type:jsonb;not null"`
	StartDate       time.Time           `gorm:"type:date;not null"`
	EndDate         *time.Time          `gorm:"type:date"`
	IsActive        bool                `gorm:"not null;index"`
	CreatedAt       time.Time           `gorm:"not null"`
	UpdatedAt       time.Time           `gorm:"not null"`
}

// TableName возвращает имя таблицы
func (medicationModel) TableName() string {
	return "medications"
}

// toDomain преобразует модель БД в доменную сущность
func (m *medicationModel) toDomain() *medication.Medication {
	var scheduleType medication.ScheduleType
	switch m.ScheduleType {
	case "daily":
		scheduleType = medication.ScheduleTypeDaily
	case "weekly":
		scheduleType = medication.ScheduleTypeWeekly
	default:
		scheduleType = medication.ScheduleTypeAsNeeded
	}

	return &medication.Medication{
		ID:              m.ID,
		UserID:          m.UserID,
		Name:            m.Name,
		Dosage:          m.Dosage,
		ScheduleType:    scheduleType,
		ScheduleDetails: medication.ScheduleDetails(m.ScheduleDetails),
		StartDate:       m.StartDate,
		EndDate:         m.EndDate,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// fromDomain преобразует доменную сущность в модель БД
func (m *medicationModel) fromDomain(med *medication.Medication) {
	m.ID = med.ID
	m.UserID = med.UserID
	m.Name = med.Name
	m.Dosage = med.Dosage
	m.ScheduleType = string(med.ScheduleType)
	m.ScheduleDetails = scheduleDetailsJSON(med.ScheduleDetails)
	m.StartDate = med.StartDate
	m.EndDate = med.EndDate
	m.IsActive = med.IsActive
	m.CreatedAt = med.CreatedAt
	m.UpdatedAt = med.UpdatedAt
}

// MedicationRepository реализует medication.Repository для PostgreSQL
type MedicationRepository struct {
	db *gorm.DB
}

// NewMedicationRepository создаёт новый репозиторий лекарств
func NewMedicationRepository(db *gorm.DB) medication.Repository {
	return &MedicationRepository{db: db}
}

// Create создаёт новое лекарство
func (r *MedicationRepository) Create(ctx context.Context, med *medication.Medication) error {
	model := &medicationModel{}
	model.fromDomain(med)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*med = *model.toDomain()
	return nil
}

// GetByID возвращает лекарство по ID
func (r *MedicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*medication.Medication, error) {
	var model medicationModel
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

// FindByUserID возвращает все лекарства пользователя
func (r *MedicationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]*medication.Medication, error) {
	query := r.db.WithContext(ctx).Model(&medicationModel{}).
		Where("user_id = ?", userID)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var models []medicationModel
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	medications := make([]*medication.Medication, len(models))
	for i := range models {
		medications[i] = models[i].toDomain()
	}

	return medications, nil
}

// Update обновляет лекарство
func (r *MedicationRepository) Update(ctx context.Context, med *medication.Medication) error {
	model := &medicationModel{}
	model.fromDomain(med)

	return r.db.WithContext(ctx).
		Model(&medicationModel{}).
		Where("id = ?", med.ID).
		Updates(model).Error
}

// Delete удаляет лекарство
func (r *MedicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&medicationModel{}).Error
}

