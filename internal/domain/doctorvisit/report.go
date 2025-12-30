package doctorvisit

import (
	"time"

	"github.com/google/uuid"
)

// Report представляет отчёт для визита к врачу
type Report struct {
	VisitID      uuid.UUID
	VisitDate    time.Time
	Period       DateRange
	Symptoms     []ReportSymptom
	WellbeingTrend WellbeingTrend
	Analyses     []ReportAnalysis
	Medications  []ReportMedication
	Questions    *string
	GeneratedAt  time.Time
}

// ReportSymptom представляет симптом в отчёте
type ReportSymptom struct {
	ID             uuid.UUID
	DateTime       time.Time
	Description    string
	WellbeingScale int
}

// ReportAnalysis представляет анализ в отчёте
type ReportAnalysis struct {
	ID        uuid.UUID
	Type      string
	Name      string
	DateTaken time.Time
}

// ReportMedication представляет лекарство в отчёте
type ReportMedication struct {
	ID       uuid.UUID
	Name     string
	Dosage   string
	IsActive bool
}

// WellbeingTrend представляет тренд самочувствия
type WellbeingTrend struct {
	Average   float64
	Min       int
	Max       int
	DataPoints []WellbeingDataPoint
}

// WellbeingDataPoint представляет точку данных для графика
type WellbeingDataPoint struct {
	Date  time.Time
	Value int
}

// NewReport создаёт новый отчёт
func NewReport(
	visitID uuid.UUID,
	visitDate time.Time,
	period DateRange,
) *Report {
	return &Report{
		VisitID:     visitID,
		VisitDate:   visitDate,
		Period:      period,
		Symptoms:    []ReportSymptom{},
		Analyses:    []ReportAnalysis{},
		Medications: []ReportMedication{},
		GeneratedAt: time.Now(),
	}
}

// AddSymptom добавляет симптом в отчёт
func (r *Report) AddSymptom(symptom ReportSymptom) {
	r.Symptoms = append(r.Symptoms, symptom)
}

// AddAnalysis добавляет анализ в отчёт
func (r *Report) AddAnalysis(analysis ReportAnalysis) {
	r.Analyses = append(r.Analyses, analysis)
}

// AddMedication добавляет лекарство в отчёт
func (r *Report) AddMedication(medication ReportMedication) {
	r.Medications = append(r.Medications, medication)
}

// SetWellbeingTrend устанавливает тренд самочувствия
func (r *Report) SetWellbeingTrend(trend WellbeingTrend) {
	r.WellbeingTrend = trend
}

// SetQuestions устанавливает вопросы к врачу
func (r *Report) SetQuestions(questions string) {
	r.Questions = &questions
}

