package usecase

import (
	"context"
	"errors"

	"github.com/kshitij-nehete/astro-report/internal/domain"
	"github.com/kshitij-nehete/astro-report/internal/repository"
)

type ReportUsecase struct {
	reportRepo repository.ReportRepository
	userRepo   repository.UserRepository
}

func NewReportUsecase(
	reportRepo repository.ReportRepository,
	userRepo repository.UserRepository,
) *ReportUsecase {
	return &ReportUsecase{
		reportRepo: reportRepo,
		userRepo:   userRepo,
	}
}

func (u *ReportUsecase) CreateReport(
	ctx context.Context,
	userID string,
	name string,
	birthDate string,
	birthTime string,
	location string,
) (*domain.Report, error) {

	// Check report count limit (max 3 for now)
	count, err := u.reportRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if count >= 3 {
		return nil, errors.New("report limit reached")
	}

	report := &domain.Report{
		UserID:    userID,
		Name:      name,
		BirthDate: birthDate,
		BirthTime: birthTime,
		Location:  location,
		Status:    domain.StatusProcessing,
	}

	err = u.reportRepo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}
