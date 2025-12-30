package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/symptom"
	"gorm.io/gorm"
)

// symptomModel представляет модель записи симптома в БД
type symptomModel struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID                uuid.UUID  `gorm:"type:uuid;not null;index"`
	DateTime              time.Time  `gorm:"not null;index"`
	Description           string     `gorm:"type:text;not null"`
	WellbeingScale        int        `gorm:"not null;check:wellbeing_scale >= 1 AND wellbeing_scale <= 10"`
	Temperature           *float64   `gorm:"type:decimal(4,1)"`
	BloodPressureSystolic *int
	BloodPressureDiastolic *int
	Pulse                 *int
	PhotoURL              *string    `gorm:"type:varchar(500)"`
	CreatedAt             time.Time  `gorm:"not null"`
	UpdatedAt             time.Time  `gorm:"not null"`
}

// TableName возвращает имя таблицы
func (symptomModel) TableName() string {
	return "symptom_entries"
}

// toDomain преобразует модель БД в доменную сущность
func (m *symptomModel) toDomain() *symptom.SymptomEntry {
	return &symptom.SymptomEntry{
		ID:                    m.ID,
		UserID:                m.UserID,
		DateTime:              m.DateTime,
		Description:           m.Description,
		WellbeingScale:        m.WellbeingScale,
		Temperature:           m.Temperature,
		BloodPressureSystolic: m.BloodPressureSystolic,
		BloodPressureDiastolic: m.BloodPressureDiastolic,
		Pulse:                 m.Pulse,
		PhotoURL:              m.PhotoURL,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

// fromDomain преобразует доменную сущность в модель БД
func (m *symptomModel) fromDomain(s *symptom.SymptomEntry) {
	m.ID = s.ID
	m.UserID = s.UserID
	m.DateTime = s.DateTime
	m.Description = s.Description
	m.WellbeingScale = s.WellbeingScale
	m.Temperature = s.Temperature
	m.BloodPressureSystolic = s.BloodPressureSystolic
	m.BloodPressureDiastolic = s.BloodPressureDiastolic
	m.Pulse = s.Pulse
	m.PhotoURL = s.PhotoURL
	m.CreatedAt = s.CreatedAt
	m.UpdatedAt = s.UpdatedAt
}

// SymptomRepository реализует symptom.Repository для PostgreSQL
type SymptomRepository struct {
	db *gorm.DB
}

// NewSymptomRepository создаёт новый репозиторий симптомов
func NewSymptomRepository(db *gorm.DB) symptom.Repository {
	return &SymptomRepository{db: db}
}

// Create создаёт новую запись симптома
func (r *SymptomRepository) Create(ctx context.Context, entry *symptom.SymptomEntry) error {
	model := &symptomModel{}
	model.fromDomain(entry)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*entry = *model.toDomain()
	return nil
}

// GetByID возвращает запись по ID
func (r *SymptomRepository) GetByID(ctx context.Context, id uuid.UUID) (*symptom.SymptomEntry, error) {
	var model symptomModel
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, symptom.ErrSymptomNotFound
		}
		return nil, err
	}

	return model.toDomain(), nil
}

// FindByFilter возвращает записи по фильтру
func (r *SymptomRepository) FindByFilter(ctx context.Context, filter symptom.Filter, limit, offset int) ([]*symptom.SymptomEntry, int, error) {
	query := r.db.WithContext(ctx).Model(&symptomModel{}).
		Where("user_id = ?", filter.UserID)

	if filter.StartDate != nil {
		query = query.Where("date_time >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date_time <= ?", *filter.EndDate)
	}
	if filter.MinWellbeingScale != nil {
		query = query.Where("wellbeing_scale >= ?", *filter.MinWellbeingScale)
	}
	if filter.MaxWellbeingScale != nil {
		query = query.Where("wellbeing_scale <= ?", *filter.MaxWellbeingScale)
	}

	// Подсчёт общего количества
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получение записей с пагинацией
	var models []symptomModel
	if err := query.
		Order("date_time DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entries := make([]*symptom.SymptomEntry, len(models))
	for i := range models {
		entries[i] = models[i].toDomain()
	}

	return entries, int(total), nil
}

// Update обновляет запись
func (r *SymptomRepository) Update(ctx context.Context, entry *symptom.SymptomEntry) error {
	model := &symptomModel{}
	model.fromDomain(entry)

	return r.db.WithContext(ctx).
		Model(&symptomModel{}).
		Where("id = ?", entry.ID).
		Updates(model).Error
}

// Delete удаляет запись
func (r *SymptomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&symptomModel{}).Error
}

// GetWellbeingTrend возвращает тренд самочувствия за период
func (r *SymptomRepository) GetWellbeingTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]symptom.WellbeingDataPoint, error) {
	var results []struct {
		Date  time.Time `gorm:"column:date"`
		Value int       `gorm:"column:value"`
	}

	err := r.db.WithContext(ctx).
		Model(&symptomModel{}).
		Select("DATE(date_time) as date, wellbeing_scale as value").
		Where("user_id = ? AND date_time >= ? AND date_time <= ?", userID, startDate, endDate).
		Order("date ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	points := make([]symptom.WellbeingDataPoint, len(results))
	for i, r := range results {
		points[i] = symptom.WellbeingDataPoint{
			Date:  r.Date,
			Value: r.Value,
		}
	}

	return points, nil
}

