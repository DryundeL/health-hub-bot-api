package graphql

import (
	"github.com/health-hub-bot-api/internal/domain/analysis"
	"github.com/health-hub-bot-api/internal/domain/doctorvisit"
	"github.com/health-hub-bot-api/internal/domain/medication"
	"github.com/health-hub-bot-api/internal/domain/symptom"
	"github.com/health-hub-bot-api/internal/domain/user"
)

// Resolver содержит зависимости для GraphQL resolvers
type Resolver struct {
	// Repositories
	userRepo        user.Repository
	symptomRepo     symptom.Repository
	analysisRepo    analysis.Repository
	medicationRepo  medication.Repository
	intakeRepo      medication.IntakeRepository
	doctorVisitRepo doctorvisit.Repository

	// Services (use cases) будут добавлены позже
}

// NewResolver создаёт новый resolver
func NewResolver(
	userRepo user.Repository,
	symptomRepo symptom.Repository,
	analysisRepo analysis.Repository,
	medicationRepo medication.Repository,
	intakeRepo medication.IntakeRepository,
	doctorVisitRepo doctorvisit.Repository,
) *Resolver {
	return &Resolver{
		userRepo:        userRepo,
		symptomRepo:     symptomRepo,
		analysisRepo:    analysisRepo,
		medicationRepo:  medicationRepo,
		intakeRepo:      intakeRepo,
		doctorVisitRepo: doctorVisitRepo,
	}
}
