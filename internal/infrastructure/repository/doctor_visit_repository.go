package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/doctorvisit"
	"gorm.io/gorm"
)

// reportDataJSON представляет JSON для report_data
type reportDataJSON struct {
	Period       dateRangeJSON   `json:"period"`
	SymptomIDs   []uuid.UUID     `json:"symptom_ids"`
	AnalysisIDs  []uuid.UUID     `json:"analysis_ids"`
	MedicationIDs []uuid.UUID    `json:"medication_ids"`
}

type dateRangeJSON struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Value реализует driver.Valuer для GORM
func (r reportDataJSON) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan реализует sql.Scanner для GORM
func (r *reportDataJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// doctorVisitModel представляет модель визита к врачу в БД
type doctorVisitModel struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID            uuid.UUID       `gorm:"type:uuid;not null;index"`
	VisitDate         time.Time       `gorm:"type:date;not null;index"`
	DoctorName        *string         `gorm:"type:varchar(255)"`
	Specialty         *string         `gorm:"type:varchar(100)"`
	Questions         *string         `gorm:"type:text"`
	ReportGeneratedAt *time.Time
	ReportData        *reportDataJSON `gorm:"type:jsonb"`
	CreatedAt         time.Time        `gorm:"not null"`
	UpdatedAt         time.Time        `gorm:"not null"`
}

// TableName возвращает имя таблицы
func (doctorVisitModel) TableName() string {
	return "doctor_visits"
}

// toDomain преобразует модель БД в доменную сущность
func (m *doctorVisitModel) toDomain() (*doctorvisit.DoctorVisit, error) {
	var reportData *doctorvisit.ReportData
	if m.ReportData != nil {
		reportData = &doctorvisit.ReportData{
			Period: doctorvisit.DateRange{
				StartDate: m.ReportData.Period.StartDate,
				EndDate:   m.ReportData.Period.EndDate,
			},
			SymptomIDs:   m.ReportData.SymptomIDs,
			AnalysisIDs:  m.ReportData.AnalysisIDs,
			MedicationIDs: m.ReportData.MedicationIDs,
		}
	}

	return &doctorvisit.DoctorVisit{
		ID:                m.ID,
		UserID:            m.UserID,
		VisitDate:         m.VisitDate,
		DoctorName:        m.DoctorName,
		Specialty:         m.Specialty,
		Questions:         m.Questions,
		ReportGeneratedAt: m.ReportGeneratedAt,
		ReportData:        reportData,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}, nil
}

// fromDomain преобразует доменную сущность в модель БД
func (m *doctorVisitModel) fromDomain(v *doctorvisit.DoctorVisit) error {
	m.ID = v.ID
	m.UserID = v.UserID
	m.VisitDate = v.VisitDate
	m.DoctorName = v.DoctorName
	m.Specialty = v.Specialty
	m.Questions = v.Questions
	m.ReportGeneratedAt = v.ReportGeneratedAt
	m.CreatedAt = v.CreatedAt
	m.UpdatedAt = v.UpdatedAt

	if v.ReportData != nil {
		m.ReportData = &reportDataJSON{
			Period: dateRangeJSON{
				StartDate: v.ReportData.Period.StartDate,
				EndDate:   v.ReportData.Period.EndDate,
			},
			SymptomIDs:   v.ReportData.SymptomIDs,
			AnalysisIDs:  v.ReportData.AnalysisIDs,
			MedicationIDs: v.ReportData.MedicationIDs,
		}
	}

	return nil
}

// DoctorVisitRepository реализует doctorvisit.Repository для PostgreSQL
type DoctorVisitRepository struct {
	db *gorm.DB
}

// NewDoctorVisitRepository создаёт новый репозиторий визитов к врачу
func NewDoctorVisitRepository(db *gorm.DB) doctorvisit.Repository {
	return &DoctorVisitRepository{db: db}
}

// Create создаёт новый визит
func (r *DoctorVisitRepository) Create(ctx context.Context, visit *doctorvisit.DoctorVisit) error {
	model := &doctorVisitModel{}
	if err := model.fromDomain(visit); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	domain, err := model.toDomain()
	if err != nil {
		return err
	}
	*visit = *domain
	return nil
}

// GetByID возвращает визит по ID
func (r *DoctorVisitRepository) GetByID(ctx context.Context, id uuid.UUID) (*doctorvisit.DoctorVisit, error) {
	var model doctorVisitModel
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, doctorvisit.ErrVisitNotFound
		}
		return nil, err
	}

	return model.toDomain()
}

// FindByUserID возвращает визиты пользователя
func (r *DoctorVisitRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*doctorvisit.DoctorVisit, int, error) {
	query := r.db.WithContext(ctx).Model(&doctorVisitModel{}).
		Where("user_id = ?", userID)

	// Подсчёт общего количества
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получение визитов с пагинацией
	var models []doctorVisitModel
	if err := query.
		Order("visit_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	visits := make([]*doctorvisit.DoctorVisit, 0, len(models))
	for i := range models {
		visit, err := models[i].toDomain()
		if err != nil {
			return nil, 0, err
		}
		visits = append(visits, visit)
	}

	return visits, int(total), nil
}

// Update обновляет визит
func (r *DoctorVisitRepository) Update(ctx context.Context, visit *doctorvisit.DoctorVisit) error {
	model := &doctorVisitModel{}
	if err := model.fromDomain(visit); err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&doctorVisitModel{}).
		Where("id = ?", visit.ID).
		Updates(model).Error
}

// Delete удаляет визит
func (r *DoctorVisitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&doctorVisitModel{}).Error
}

// GetUpcomingVisits возвращает предстоящие визиты
func (r *DoctorVisitRepository) GetUpcomingVisits(ctx context.Context, userID uuid.UUID, beforeDate time.Time) ([]*doctorvisit.DoctorVisit, error) {
	var models []doctorVisitModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND visit_date >= ? AND visit_date <= ?", userID, time.Now(), beforeDate).
		Order("visit_date ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	visits := make([]*doctorvisit.DoctorVisit, 0, len(models))
	for i := range models {
		visit, err := models[i].toDomain()
		if err != nil {
			return nil, err
		}
		visits = append(visits, visit)
	}

	return visits, nil
}

