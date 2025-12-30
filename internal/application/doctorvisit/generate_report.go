package doctorvisit

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/health-hub-bot-api/internal/domain/analysis"
	"github.com/health-hub-bot-api/internal/domain/doctorvisit"
	"github.com/health-hub-bot-api/internal/domain/medication"
	"github.com/health-hub-bot-api/internal/domain/symptom"
)

// GenerateReportUseCase представляет use case для генерации отчёта к врачу
type GenerateReportUseCase struct {
	doctorVisitRepo doctorvisit.Repository
	symptomRepo     symptom.Repository
	analysisRepo    analysis.Repository
	medicationRepo  medication.Repository
}

// NewGenerateReportUseCase создаёт новый use case
func NewGenerateReportUseCase(
	doctorVisitRepo doctorvisit.Repository,
	symptomRepo symptom.Repository,
	analysisRepo analysis.Repository,
	medicationRepo medication.Repository,
) *GenerateReportUseCase {
	return &GenerateReportUseCase{
		doctorVisitRepo: doctorVisitRepo,
		symptomRepo:     symptomRepo,
		analysisRepo:    analysisRepo,
		medicationRepo:  medicationRepo,
	}
}

// GenerateReportInput представляет входные данные для генерации отчёта
type GenerateReportInput struct {
	VisitID   uuid.UUID
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	Questions *string
}

// Execute выполняет генерацию отчёта
func (uc *GenerateReportUseCase) Execute(ctx context.Context, input GenerateReportInput) (*doctorvisit.Report, error) {
	// Получаем визит
	visit, err := uc.doctorVisitRepo.GetByID(ctx, input.VisitID)
	if err != nil {
		return nil, err
	}

	// Проверяем, что визит принадлежит пользователю
	if visit.UserID != input.UserID {
		return nil, ErrUnauthorized
	}

	// Создаём отчёт
	report := doctorvisit.NewReport(
		visit.ID,
		visit.VisitDate,
		doctorvisit.DateRange{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
		},
	)

	// Получаем симптомы за период
	symptomFilter := symptom.Filter{
		UserID:    input.UserID,
		StartDate: &input.StartDate,
		EndDate:   &input.EndDate,
	}
	symptoms, _, err := uc.symptomRepo.FindByFilter(ctx, symptomFilter, 1000, 0)
	if err != nil {
		return nil, err
	}

	for _, s := range symptoms {
		report.AddSymptom(doctorvisit.ReportSymptom{
			ID:             s.ID,
			DateTime:       s.DateTime,
			Description:    s.Description,
			WellbeingScale: s.WellbeingScale,
		})
	}

	// Получаем тренд самочувствия
	trendData, err := uc.symptomRepo.GetWellbeingTrend(ctx, input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	// Вычисляем статистику
	if len(trendData) > 0 {
		var sum int
		min := trendData[0].Value
		max := trendData[0].Value
		dataPoints := make([]doctorvisit.WellbeingDataPoint, 0, len(trendData))

		for _, point := range trendData {
			sum += point.Value
			if point.Value < min {
				min = point.Value
			}
			if point.Value > max {
				max = point.Value
			}
			dataPoints = append(dataPoints, doctorvisit.WellbeingDataPoint{
				Date:  point.Date,
				Value: point.Value,
			})
		}

		average := float64(sum) / float64(len(trendData))
		report.SetWellbeingTrend(doctorvisit.WellbeingTrend{
			Average:    average,
			Min:        min,
			Max:        max,
			DataPoints: dataPoints,
		})
	}

	// Получаем анализы за период
	analysisFilter := analysis.Filter{
		UserID:    input.UserID,
		StartDate: &input.StartDate,
		EndDate:   &input.EndDate,
	}
	analyses, _, err := uc.analysisRepo.FindByFilter(ctx, analysisFilter, 1000, 0)
	if err != nil {
		return nil, err
	}

	for _, a := range analyses {
		report.AddAnalysis(doctorvisit.ReportAnalysis{
			ID:        a.ID,
			Type:      string(a.Type),
			Name:      a.Name,
			DateTaken: a.DateTaken,
		})
	}

	// Получаем активные лекарства
	medications, err := uc.medicationRepo.FindByUserID(ctx, input.UserID, true)
	if err != nil {
		return nil, err
	}

	for _, m := range medications {
		report.AddMedication(doctorvisit.ReportMedication{
			ID:       m.ID,
			Name:     m.Name,
			Dosage:   m.Dosage,
			IsActive: m.IsActive,
		})
	}

	// Устанавливаем вопросы, если есть
	if input.Questions != nil {
		report.SetQuestions(*input.Questions)
		visit.Questions = input.Questions
	}

	// Сохраняем данные отчёта в визит
	reportData := doctorvisit.ReportData{
		Period: doctorvisit.DateRange{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
		},
		SymptomIDs:   make([]uuid.UUID, 0, len(symptoms)),
		AnalysisIDs:  make([]uuid.UUID, 0, len(analyses)),
		MedicationIDs: make([]uuid.UUID, 0, len(medications)),
	}

	for _, s := range symptoms {
		reportData.SymptomIDs = append(reportData.SymptomIDs, s.ID)
	}
	for _, a := range analyses {
		reportData.AnalysisIDs = append(reportData.AnalysisIDs, a.ID)
	}
	for _, m := range medications {
		reportData.MedicationIDs = append(reportData.MedicationIDs, m.ID)
	}

	if err := visit.SetReportData(reportData); err != nil {
		return nil, err
	}

	// Обновляем визит
	if err := uc.doctorVisitRepo.Update(ctx, visit); err != nil {
		return nil, err
	}

	return report, nil
}

var ErrUnauthorized = doctorvisit.ErrUnauthorized

