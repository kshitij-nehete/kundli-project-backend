package repository

import (
	"context"

	"github.com/kshitij-nehete/astro-report/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
