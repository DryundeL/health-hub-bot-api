package doctorvisit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DoctorVisit представляет визит к врачу
type DoctorVisit struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	VisitDate         time.Time
	DoctorName        *string
	Specialty         *string
	Questions         *string
	ReportGeneratedAt *time.Time
	ReportData        *ReportData
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ReportData представляет данные отчёта для визита
type ReportData struct {
	Period      DateRange
	SymptomIDs  []uuid.UUID
	AnalysisIDs []uuid.UUID
	MedicationIDs []uuid.UUID
}

// DateRange представляет диапазон дат
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// NewDoctorVisit создаёт новый визит к врачу
func NewDoctorVisit(
	userID uuid.UUID,
	visitDate time.Time,
) *DoctorVisit {
	now := time.Now()
	return &DoctorVisit{
		ID:        uuid.New(),
		UserID:    userID,
		VisitDate: visitDate,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update обновляет визит
func (d *DoctorVisit) Update(
	visitDate *time.Time,
	doctorName *string,
	specialty *string,
	questions *string,
) {
	if visitDate != nil {
		d.VisitDate = *visitDate
	}
	if doctorName != nil {
		d.DoctorName = doctorName
	}
	if specialty != nil {
		d.Specialty = specialty
	}
	if questions != nil {
		d.Questions = questions
	}
	d.UpdatedAt = time.Now()
}

// SetReportData сохраняет данные отчёта
func (d *DoctorVisit) SetReportData(data ReportData) error {
	d.ReportData = &data
	now := time.Now()
	d.ReportGeneratedAt = &now
	d.UpdatedAt = now
	return nil
}

// GetReportDataJSON возвращает данные отчёта в формате JSON
func (d *DoctorVisit) GetReportDataJSON() ([]byte, error) {
	if d.ReportData == nil {
		return nil, nil
	}
	return json.Marshal(d.ReportData)
}

