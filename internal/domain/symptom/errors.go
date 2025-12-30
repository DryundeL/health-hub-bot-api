package symptom

import "errors"

var (
	ErrInvalidWellbeingScale = errors.New("wellbeing scale must be between 1 and 10")
	ErrSymptomNotFound       = errors.New("symptom entry not found")
)

