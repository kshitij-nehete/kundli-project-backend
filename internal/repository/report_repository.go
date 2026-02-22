package repository

import (
	"context"

	"github.com/kshitij-nehete/astro-report/internal/domain"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	CountByUser(ctx context.Context, userID string) (int64, error)
}
