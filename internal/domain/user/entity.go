package user

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя системы
type User struct {
	ID             uuid.UUID
	TelegramUserID int64
	Name           string
	Age            *int
	Gender         *Gender
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// Gender представляет пол пользователя
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// NewUser создаёт нового пользователя
func NewUser(telegramUserID int64, name string) *User {
	now := time.Now()
	return &User{
		ID:             uuid.New(),
		TelegramUserID: telegramUserID,
		Name:           name,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// UpdateProfile обновляет профиль пользователя
func (u *User) UpdateProfile(name string, age *int, gender *Gender) {
	if name != "" {
		u.Name = name
	}
	if age != nil {
		u.Age = age
	}
	if gender != nil {
		u.Gender = gender
	}
	u.UpdatedAt = time.Now()
}

// IsDeleted проверяет, удалён ли пользователь
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// Delete помечает пользователя как удалённого
func (u *User) Delete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}
