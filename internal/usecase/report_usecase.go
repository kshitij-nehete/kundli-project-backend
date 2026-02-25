package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kshitij-nehete/astro-report/internal/ai"
	"github.com/kshitij-nehete/astro-report/internal/domain"
	"github.com/kshitij-nehete/astro-report/internal/repository"
)

type ReportUsecase struct {
	reportRepo   repository.ReportRepository
	userRepo     repository.UserRepository
	orchestrator *ai.Orchestrator
}

func NewReportUsecase(
	reportRepo repository.ReportRepository,
	userRepo repository.UserRepository,
	orchestrator *ai.Orchestrator,
) *ReportUsecase {
	return &ReportUsecase{
		reportRepo:   reportRepo,
		userRepo:     userRepo,
		orchestrator: orchestrator,
	}
}

func (u *ReportUsecase) CreateReport(
	ctx context.Context,
	userID string,
	name string,
	birthDate string,
	birthTime string,
	placeOfBirth string,
) (*domain.Report, error) {

	count, err := u.reportRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if count >= 3 {
		return nil, errors.New("report limit reached")
	}

	input := map[string]interface{}{
		"name":       name,
		"birth_date": birthDate,
	}

	if birthTime != "" {
		input["birth_time"] = birthTime
	}
	if placeOfBirth != "" {
		input["place_of_birth"] = placeOfBirth
	}

	// 🔥 Context timeout enforcement
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	// 🔥 Retry wrapper
	resultMap, err := ai.WithRetry(ctxWithTimeout, 2, func() (map[string]interface{}, error) {
		return u.orchestrator.Run(ctxWithTimeout, input)
	})

	if err != nil {
		return nil, errors.New("failed to generate report")
	}

	// 🔥 Strict JSON validation
	var structured domain.NumerologyReport
	bytes, _ := json.Marshal(resultMap)

	if err := json.Unmarshal(bytes, &structured); err != nil {
		return nil, errors.New("invalid report structure from AI")
	}

	report := &domain.Report{
		UserID:         userID,
		Name:           name,
		BirthDate:      birthDate,
		BirthTime:      birthTime,
		PlaceOfBirth:   placeOfBirth,
		Status:         domain.StatusActive,
		NumerologyData: &structured,
	}

	err = u.reportRepo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}
