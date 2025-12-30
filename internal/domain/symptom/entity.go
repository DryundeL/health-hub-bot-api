package symptom

import (
	"time"

	"github.com/google/uuid"
)

// SymptomEntry представляет запись симптома в дневнике
type SymptomEntry struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	DateTime              time.Time
	Description           string
	WellbeingScale        int // 1-10
	Temperature           *float64
	BloodPressureSystolic *int
	BloodPressureDiastolic *int
	Pulse                 *int
	PhotoURL              *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// NewSymptomEntry создаёт новую запись симптома
func NewSymptomEntry(
	userID uuid.UUID,
	dateTime time.Time,
	description string,
	wellbeingScale int,
) (*SymptomEntry, error) {
	if err := validateWellbeingScale(wellbeingScale); err != nil {
		return nil, err
	}

	now := time.Now()
	return &SymptomEntry{
		ID:             uuid.New(),
		UserID:         userID,
		DateTime:       dateTime,
		Description:    description,
		WellbeingScale: wellbeingScale,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Update обновляет запись симптома
func (s *SymptomEntry) Update(
	dateTime *time.Time,
	description *string,
	wellbeingScale *int,
	temperature *float64,
	bloodPressureSystolic *int,
	bloodPressureDiastolic *int,
	pulse *int,
	photoURL *string,
) error {
	if dateTime != nil {
		s.DateTime = *dateTime
	}
	if description != nil {
		s.Description = *description
	}
	if wellbeingScale != nil {
		if err := validateWellbeingScale(*wellbeingScale); err != nil {
			return err
		}
		s.WellbeingScale = *wellbeingScale
	}
	if temperature != nil {
		s.Temperature = temperature
	}
	if bloodPressureSystolic != nil {
		s.BloodPressureSystolic = bloodPressureSystolic
	}
	if bloodPressureDiastolic != nil {
		s.BloodPressureDiastolic = bloodPressureDiastolic
	}
	if pulse != nil {
		s.Pulse = pulse
	}
	if photoURL != nil {
		s.PhotoURL = photoURL
	}
	s.UpdatedAt = time.Now()
	return nil
}

// validateWellbeingScale проверяет корректность шкалы самочувствия
func validateWellbeingScale(scale int) error {
	if scale < 1 || scale > 10 {
		return ErrInvalidWellbeingScale
	}
	return nil
}

