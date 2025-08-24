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
	"os"
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

	migrationPath := os.Getenv("MIGRATIONS_PATH")
	if migrationPath == "" {
		migrationPath = "./migrations"
	}
	var m *migrate.Migrate

	for i := 0; i < 10; i++ {
		m, err = migrate.New(migrationPath, connString)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logr.Info("migrated successfully", zap.String("path", migrationPath))

	return &Storage{db: dbPool}, nil
}

func (s *Storage) SaveUserDB(ctx context.Context, email string, passHash []byte, role string) (int64, error) {

	var userID int64
	err := s.db.QueryRow(ctx, `INSERT INTO users(email, password, role) VALUES($1, $2, $3) RETURNING id`, email, passHash, role).Scan(&userID)
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

	err := s.db.QueryRow(ctx, `SELECT id, email, password, role FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PassHash, &user.Role)
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
