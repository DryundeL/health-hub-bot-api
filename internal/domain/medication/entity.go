package medication

import (
	"time"

	"github.com/google/uuid"
)

// Medication представляет лекарство
type Medication struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Name          string
	Dosage        string
	ScheduleType  ScheduleType
	ScheduleDetails ScheduleDetails
	StartDate     time.Time
	EndDate       *time.Time
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ScheduleType представляет тип расписания приёма
type ScheduleType string

const (
	ScheduleTypeDaily     ScheduleType = "daily"
	ScheduleTypeWeekly    ScheduleType = "weekly"
	ScheduleTypeAsNeeded  ScheduleType = "as_needed"
)

// ScheduleDetails представляет детали расписания
type ScheduleDetails struct {
	Times []string // ["09:00", "21:00"]
	Days  []int    // [1,2,3,4,5,6,7] для дней недели (1=Monday)
}

// NewMedication создаёт новое лекарство
func NewMedication(
	userID uuid.UUID,
	name string,
	dosage string,
	scheduleType ScheduleType,
	scheduleDetails ScheduleDetails,
	startDate time.Time,
) *Medication {
	now := time.Now()
	return &Medication{
		ID:              uuid.New(),
		UserID:          userID,
		Name:            name,
		Dosage:          dosage,
		ScheduleType:    scheduleType,
		ScheduleDetails: scheduleDetails,
		StartDate:       startDate,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Update обновляет лекарство
func (m *Medication) Update(
	name *string,
	dosage *string,
	scheduleType *ScheduleType,
	scheduleDetails *ScheduleDetails,
	startDate *time.Time,
	endDate *time.Time,
	isActive *bool,
) {
	if name != nil {
		m.Name = *name
	}
	if dosage != nil {
		m.Dosage = *dosage
	}
	if scheduleType != nil {
		m.ScheduleType = *scheduleType
	}
	if scheduleDetails != nil {
		m.ScheduleDetails = *scheduleDetails
	}
	if startDate != nil {
		m.StartDate = *startDate
	}
	if endDate != nil {
		m.EndDate = endDate
	}
	if isActive != nil {
		m.IsActive = *isActive
	}
	m.UpdatedAt = time.Now()
}

// Deactivate деактивирует лекарство
func (m *Medication) Deactivate() {
	m.IsActive = false
	m.UpdatedAt = time.Now()
}

// IsExpired проверяет, истёк ли срок приёма лекарства
func (m *Medication) IsExpired() bool {
	if m.EndDate == nil {
		return false
	}
	return time.Now().After(*m.EndDate)
}

