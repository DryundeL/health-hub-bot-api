package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/analysis"
	"gorm.io/gorm"
)

// analysisModel представляет модель анализа в БД
type analysisModel struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	Type            string     `gorm:"type:varchar(20);not null;check:type IN ('blood','urine','ultrasound','xray','other');index"`
	Name            string     `gorm:"type:varchar(255);not null"`
	DateTaken       time.Time  `gorm:"type:date;not null;index"`
	FileURL         string     `gorm:"type:varchar(500);not null"`
	FileType        string     `gorm:"type:varchar(10);not null;check:file_type IN ('image','pdf')"`
	NextReminderDate *time.Time `gorm:"type:date;index"`
	CreatedAt       time.Time  `gorm:"not null"`
	UpdatedAt       time.Time  `gorm:"not null"`
}

// TableName возвращает имя таблицы
func (analysisModel) TableName() string {
	return "analyses"
}

// toDomain преобразует модель БД в доменную сущность
func (m *analysisModel) toDomain() *analysis.Analysis {
	var analysisType analysis.Type
	switch m.Type {
	case "blood":
		analysisType = analysis.TypeBlood
	case "urine":
		analysisType = analysis.TypeUrine
	case "ultrasound":
		analysisType = analysis.TypeUltrasound
	case "xray":
		analysisType = analysis.TypeXRay
	default:
		analysisType = analysis.TypeOther
	}

	var fileType analysis.FileType
	if m.FileType == "pdf" {
		fileType = analysis.FileTypePDF
	} else {
		fileType = analysis.FileTypeImage
	}

	return &analysis.Analysis{
		ID:              m.ID,
		UserID:          m.UserID,
		Type:            analysisType,
		Name:            m.Name,
		DateTaken:       m.DateTaken,
		FileURL:         m.FileURL,
		FileType:        fileType,
		NextReminderDate: m.NextReminderDate,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// fromDomain преобразует доменную сущность в модель БД
func (m *analysisModel) fromDomain(a *analysis.Analysis) {
	m.ID = a.ID
	m.UserID = a.UserID
	m.Type = string(a.Type)
	m.Name = a.Name
	m.DateTaken = a.DateTaken
	m.FileURL = a.FileURL
	m.FileType = string(a.FileType)
	m.NextReminderDate = a.NextReminderDate
	m.CreatedAt = a.CreatedAt
	m.UpdatedAt = a.UpdatedAt
}

// AnalysisRepository реализует analysis.Repository для PostgreSQL
type AnalysisRepository struct {
	db *gorm.DB
}

// NewAnalysisRepository создаёт новый репозиторий анализов
func NewAnalysisRepository(db *gorm.DB) analysis.Repository {
	return &AnalysisRepository{db: db}
}

// Create создаёт новый анализ
func (r *AnalysisRepository) Create(ctx context.Context, a *analysis.Analysis) error {
	model := &analysisModel{}
	model.fromDomain(a)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*a = *model.toDomain()
	return nil
}

// GetByID возвращает анализ по ID
func (r *AnalysisRepository) GetByID(ctx context.Context, id uuid.UUID) (*analysis.Analysis, error) {
	var model analysisModel
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

// FindByFilter возвращает анализы по фильтру
func (r *AnalysisRepository) FindByFilter(ctx context.Context, filter analysis.Filter, limit, offset int) ([]*analysis.Analysis, int, error) {
	query := r.db.WithContext(ctx).Model(&analysisModel{}).
		Where("user_id = ?", filter.UserID)

	if filter.Type != nil {
		query = query.Where("type = ?", string(*filter.Type))
	}
	if filter.StartDate != nil {
		query = query.Where("date_taken >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date_taken <= ?", *filter.EndDate)
	}

	// Подсчёт общего количества
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получение анализов с пагинацией
	var models []analysisModel
	if err := query.
		Order("date_taken DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	analyses := make([]*analysis.Analysis, len(models))
	for i := range models {
		analyses[i] = models[i].toDomain()
	}

	return analyses, int(total), nil
}

// Update обновляет анализ
func (r *AnalysisRepository) Update(ctx context.Context, a *analysis.Analysis) error {
	model := &analysisModel{}
	model.fromDomain(a)

	return r.db.WithContext(ctx).
		Model(&analysisModel{}).
		Where("id = ?", a.ID).
		Updates(model).Error
}

// Delete удаляет анализ
func (r *AnalysisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&analysisModel{}).Error
}

// GetByType группирует анализы по типу
func (r *AnalysisRepository) GetByType(ctx context.Context, userID uuid.UUID) (map[analysis.Type][]*analysis.Analysis, error) {
	var models []analysisModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date_taken DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	result := make(map[analysis.Type][]*analysis.Analysis)
	for i := range models {
		domain := models[i].toDomain()
		result[domain.Type] = append(result[domain.Type], domain)
	}

	return result, nil
}

// GetUpcomingReminders возвращает анализы с предстоящими напоминаниями
func (r *AnalysisRepository) GetUpcomingReminders(ctx context.Context, userID uuid.UUID, beforeDate time.Time) ([]*analysis.Analysis, error) {
	var models []analysisModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND next_reminder_date IS NOT NULL AND next_reminder_date <= ?", userID, beforeDate).
		Order("next_reminder_date ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	analyses := make([]*analysis.Analysis, len(models))
	for i := range models {
		analyses[i] = models[i].toDomain()
	}

	return analyses, nil
}

