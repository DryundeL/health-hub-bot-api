package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/user"
	"gorm.io/gorm"
)

// userModel представляет модель пользователя в БД
type userModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TelegramUserID int64      `gorm:"uniqueIndex;not null"`
	Name          string     `gorm:"not null"`
	Age           *int
	Gender        *string    `gorm:"type:varchar(10);check:gender IN ('male','female','other')"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     time.Time  `gorm:"not null"`
	DeletedAt     *time.Time `gorm:"index"`
}

// TableName возвращает имя таблицы
func (userModel) TableName() string {
	return "users"
}

// toDomain преобразует модель БД в доменную сущность
func (m *userModel) toDomain() *user.User {
	var gender *user.Gender
	if m.Gender != nil {
		g := user.Gender(*m.Gender)
		gender = &g
	}

	return &user.User{
		ID:            m.ID,
		TelegramUserID: m.TelegramUserID,
		Name:          m.Name,
		Age:           m.Age,
		Gender:        gender,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

// fromDomain преобразует доменную сущность в модель БД
func (m *userModel) fromDomain(u *user.User) {
	m.ID = u.ID
	m.TelegramUserID = u.TelegramUserID
	m.Name = u.Name
	m.Age = u.Age
	if u.Gender != nil {
		gender := string(*u.Gender)
		m.Gender = &gender
	}
	m.CreatedAt = u.CreatedAt
	m.UpdatedAt = u.UpdatedAt
	m.DeletedAt = u.DeletedAt
}

// UserRepository реализует user.Repository для PostgreSQL
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создаёт новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

// Create создаёт нового пользователя
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	model := &userModel{}
	model.fromDomain(u)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*u = *model.toDomain()
	return nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return model.toDomain(), nil
}

// GetByTelegramUserID возвращает пользователя по Telegram User ID
func (r *UserRepository) GetByTelegramUserID(ctx context.Context, telegramUserID int64) (*user.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).
		Where("telegram_user_id = ? AND deleted_at IS NULL", telegramUserID).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return model.toDomain(), nil
}

// Update обновляет пользователя
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	model := &userModel{}
	model.fromDomain(u)

	return r.db.WithContext(ctx).
		Model(&userModel{}).
		Where("id = ? AND deleted_at IS NULL", u.ID).
		Updates(model).Error
}

// Delete удаляет пользователя (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&userModel{}).
		Where("id = ?", id).
		Update("deleted_at", &now).Error
}

