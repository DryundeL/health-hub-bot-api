package symptom

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/symptom"
)

// CreateSymptomUseCase представляет use case для создания записи симптома
type CreateSymptomUseCase struct {
	symptomRepo symptom.Repository
	// fileStorage для загрузки фото (будет добавлен позже)
}

// NewCreateSymptomUseCase создаёт новый use case
func NewCreateSymptomUseCase(symptomRepo symptom.Repository) *CreateSymptomUseCase {
	return &CreateSymptomUseCase{
		symptomRepo: symptomRepo,
	}
}

// CreateSymptomInput представляет входные данные для создания симптома
type CreateSymptomInput struct {
	UserID                 uuid.UUID
	DateTime               time.Time
	Description            string
	WellbeingScale         int
	Temperature            *float64
	BloodPressureSystolic  *int
	BloodPressureDiastolic *int
	Pulse                  *int
	PhotoData              []byte // Будет обработан и сохранён
}

// Execute выполняет создание записи симптома
func (uc *CreateSymptomUseCase) Execute(ctx context.Context, input CreateSymptomInput) (*symptom.SymptomEntry, error) {
	// Создаём сущность
	entry, err := symptom.NewSymptomEntry(
		input.UserID,
		input.DateTime,
		input.Description,
		input.WellbeingScale,
	)
	if err != nil {
		return nil, err
	}

	// Обновляем дополнительные поля
	if err := entry.Update(
		nil, // dateTime уже установлен
		nil, // description уже установлен
		nil, // wellbeingScale уже установлен
		input.Temperature,
		input.BloodPressureSystolic,
		input.BloodPressureDiastolic,
		input.Pulse,
		nil, // photoURL будет установлен после сохранения файла
	); err != nil {
		return nil, err
	}

	// TODO: Сохранить фото, если есть
	// if len(input.PhotoData) > 0 {
	//     photoURL, err := uc.fileStorage.Save(ctx, input.PhotoData)
	//     if err != nil {
	//         return nil, err
	//     }
	//     entry.PhotoURL = &photoURL
	// }

	// Сохраняем в репозиторий
	if err := uc.symptomRepo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}
