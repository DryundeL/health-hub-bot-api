package doctorvisit

import "errors"

var (
	ErrVisitNotFound = errors.New("doctor visit not found")
	ErrUnauthorized = errors.New("unauthorized access to doctor visit")
)

