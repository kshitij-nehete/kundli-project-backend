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

	var structured domain.NumerologyReport
	bytes, _ := json.Marshal(resultMap)

	if err := json.Unmarshal(bytes, &structured); err != nil {
		return nil, errors.New("invalid report structure from AI")
	}

	// 🔥 Deterministic validation
	if err := ai.ValidateNumerologyReport(structured); err != nil {

		// 🔥 Try one regeneration
		resultMap, err = u.orchestrator.Run(ctxWithTimeout, input)
		if err != nil {
			return nil, errors.New("AI regeneration failed")
		}

		bytes, _ = json.Marshal(resultMap)
		if err := json.Unmarshal(bytes, &structured); err != nil {
			return nil, errors.New("invalid regenerated report")
		}

		if err := ai.ValidateNumerologyReport(structured); err != nil {
			return nil, errors.New("report validation failed after regeneration")
		}
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

func (u *ReportUsecase) GetReportByID(
	ctx context.Context,
	userID string,
	reportID string,
) (*domain.Report, error) {

	report, err := u.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	if report.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// 🔥 Expiry enforcement
	if report.IsExpired() && report.Status != domain.StatusExpired {

		err := u.reportRepo.UpdateStatus(
			ctx,
			report.ID,
			domain.StatusExpired,
		)
		if err == nil {
			report.Status = domain.StatusExpired
		}
	}

	return report, nil
}

func (u *ReportUsecase) GetUserReports(
	ctx context.Context,
	userID string,
) ([]*domain.Report, error) {

	reports, err := u.reportRepo.FindByUser(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	for _, report := range reports {

		if report.IsExpired() && report.Status != domain.StatusExpired {

			err := u.reportRepo.UpdateStatus(
				ctx,
				report.ID,
				domain.StatusExpired,
			)
			if err == nil {
				report.Status = domain.StatusExpired
			}
		}
	}

	return reports, nil
}
