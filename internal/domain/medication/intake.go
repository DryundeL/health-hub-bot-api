package medication

import (
	"time"

	"github.com/google/uuid"
)

// MedicationIntake представляет факт приёма лекарства
type MedicationIntake struct {
	ID            uuid.UUID
	MedicationID  uuid.UUID
	ScheduledTime time.Time
	TakenAt       *time.Time
	IsTaken       bool
	Notes         *string
	CreatedAt     time.Time
}

// NewMedicationIntake создаёт новую запись о приёме лекарства
func NewMedicationIntake(
	medicationID uuid.UUID,
	scheduledTime time.Time,
) *MedicationIntake {
	now := time.Now()
	return &MedicationIntake{
		ID:            uuid.New(),
		MedicationID:  medicationID,
		ScheduledTime: scheduledTime,
		IsTaken:       false,
		CreatedAt:     now,
	}
}

// MarkTaken отмечает приём лекарства как выполненный
func (m *MedicationIntake) MarkTaken(notes *string) {
	now := time.Now()
	m.IsTaken = true
	m.TakenAt = &now
	if notes != nil {
		m.Notes = notes
	}
}

// MarkNotTaken отмечает приём лекарства как невыполненный
func (m *MedicationIntake) MarkNotTaken() {
	m.IsTaken = false
	m.TakenAt = nil
	m.Notes = nil
}

