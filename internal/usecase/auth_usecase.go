package usecase

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kshitij-nehete/astro-report/internal/domain"
	"github.com/kshitij-nehete/astro-report/internal/repository"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}

func (u *AuthUsecase) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) error {

	// Check if user already exists
	existing, _ := u.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsPremium:    false,
		ReportCount:  0,
		CreatedAt:    time.Now(),
	}

	return u.userRepo.Create(ctx, user)
}
