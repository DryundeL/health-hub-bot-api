package analysis

import (
	"time"

	"github.com/google/uuid"
)

// Analysis представляет медицинский анализ
type Analysis struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Type            Type
	Name            string
	DateTaken       time.Time
	FileURL         string
	FileType        FileType
	NextReminderDate *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Type представляет тип анализа
type Type string

const (
	TypeBlood      Type = "blood"
	TypeUrine      Type = "urine"
	TypeUltrasound Type = "ultrasound"
	TypeXRay       Type = "xray"
	TypeOther      Type = "other"
)

// FileType представляет тип файла
type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypePDF   FileType = "pdf"
)

// NewAnalysis создаёт новый анализ
func NewAnalysis(
	userID uuid.UUID,
	analysisType Type,
	name string,
	dateTaken time.Time,
	fileURL string,
	fileType FileType,
) *Analysis {
	now := time.Now()
	return &Analysis{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      analysisType,
		Name:      name,
		DateTaken: dateTaken,
		FileURL:   fileURL,
		FileType:  fileType,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update обновляет анализ
func (a *Analysis) Update(
	analysisType *Type,
	name *string,
	dateTaken *time.Time,
	fileURL *string,
	fileType *FileType,
	nextReminderDate *time.Time,
) {
	if analysisType != nil {
		a.Type = *analysisType
	}
	if name != nil {
		a.Name = *name
	}
	if dateTaken != nil {
		a.DateTaken = *dateTaken
	}
	if fileURL != nil {
		a.FileURL = *fileURL
	}
	if fileType != nil {
		a.FileType = *fileType
	}
	if nextReminderDate != nil {
		a.NextReminderDate = nextReminderDate
	}
	a.UpdatedAt = time.Now()
}

// SetReminder устанавливает напоминание о следующем анализе
func (a *Analysis) SetReminder(date time.Time) {
	a.NextReminderDate = &date
	a.UpdatedAt = time.Now()
}

// ClearReminder удаляет напоминание
func (a *Analysis) ClearReminder() {
	a.NextReminderDate = nil
	a.UpdatedAt = time.Now()
}

