package repository

import (
	"context"

	"github.com/kshitij-nehete/astro-report/internal/domain"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	CountByUser(ctx context.Context, userID string) (int64, error)
	UpdateStatus(ctx context.Context, reportID string, status domain.ReportStatus) error
	FindByID(ctx context.Context, reportID string) (*domain.Report, error)
	FindByUser(ctx context.Context, userID string, limit int64) ([]*domain.Report, error)
}
