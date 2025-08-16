package sso

import (
	"context"
	"fmt"
	"github.com/weeweeshka/sso_for_tataisk/internal/domain/models"
	"github.com/weeweeshka/sso_for_tataisk/pkg/libs/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Sso struct {
	logr     *zap.Logger
	repo     SsoRepo
	tokenTTl time.Duration
}

type SsoRepo interface {
	UserDB(
		ctx context.Context,
		email string) (models.User, error)

	AppDB(
		ctx context.Context,
		appID int32) (models.App, error)

	SaveUserDB(ctx context.Context, email string, passHash []byte) (int64, error)

	SaveAppDB(ctx context.Context, name string, secret string) (int32, error)
}

func NewSsoService(logr *zap.Logger, repo SsoRepo, tokenTTl time.Duration) *Sso {

	return &Sso{
		logr:     logr,
		repo:     repo,
		tokenTTl: tokenTTl,
	}
}

func (s *Sso) Register(ctx context.Context, email string, password string) (int64, error) {

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	userID, err := s.repo.SaveUserDB(ctx, email, passHash)
	if err != nil {
		return 0, fmt.Errorf("s.repo.SaveUserDB: %w", err)
	}

	return userID, nil
}

func (s *Sso) Regapp(ctx context.Context, name string, secret string) (int32, error) {

	appID, err := s.repo.SaveAppDB(ctx, name, secret)
	if err != nil {
		return 0, fmt.Errorf("s.repo.SaveAppDB: %w", err)
	}

	return appID, nil

}

func (s *Sso) Login(ctx context.Context, email string, password string, appID int32) (string, error) {

	app, err := s.repo.AppDB(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("s.repo.AppDB: %w", err)
	}

	user, err := s.repo.UserDB(ctx, email)
	if err != nil {
		return "", fmt.Errorf("s.repo.UserDB: %w", err)
	}

	if bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)) != nil {
		return "", fmt.Errorf("s.repo.UserDB: invalid password")
	}

	token, err := jwt.JwtToken(user, app, s.tokenTTl)
	if err != nil {
		return "", fmt.Errorf("jwt.JwtToken: %w", err)
	}

	return token, nil

}
