package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weeweeshka/sso_for_tataisk/internal/domain/models"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

func configurationPool(config *pgxpool.Config) {
	config.MaxConns = int32(20)
	config.MinConns = int32(5)
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second
}

func NewStorage(connString string, logr *zap.Logger) (*Storage, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	configurationPool(config)

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}
	logr.Info("connected to database", zap.String("path to db", connString))

	absPath, err := filepath.Abs("../../migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to determine migration URL: %w", err)
	}

	migrationUrl := "file://" + filepath.ToSlash(absPath)

	m, err := migrate.New(migrationUrl, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	logr.Info("migrated successfully")

	return &Storage{db: dbPool}, nil
}

func (s *Storage) SaveUserDB(ctx context.Context, email string, passHash []byte) (int64, error) {

	var userID int64
	err := s.db.QueryRow(ctx, `INSERT INTO users(email, password) VALUES($1, $2) RETURNING id`, email, passHash).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}

func (s *Storage) SaveAppDB(ctx context.Context, name string, secret string) (int32, error) {

	var appID int32
	err := s.db.QueryRow(ctx, `INSERT INTO apps(name, secret) VALUES($1, $2) RETURNING id`, name, secret).Scan(&appID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert app: %w", err)
	}
	return appID, nil
}

func (s *Storage) UserDB(ctx context.Context, email string) (models.User, error) {

	var user models.User

	err := s.db.QueryRow(ctx, `SELECT id, email, password FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	return user, nil
}

func (s *Storage) AppDB(ctx context.Context, appID int32) (models.App, error) {

	var app models.App

	err := s.db.QueryRow(ctx, `SELECT id, name, secret FROM apps WHERE id = $1`, appID).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		return models.App{}, fmt.Errorf("failed to query app: %w", err)
	}

	return app, nil
}
