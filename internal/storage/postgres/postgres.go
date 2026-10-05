package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	timeout time.Duration
}

func New(ctx context.Context,
	host, port, user, password, db string,
	timeout time.Duration,
	) (*Storage, error) {
		const op = "storage.postgres.New"

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)

		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, fmt.Errorf("%s: ping: %w", op, err)
		}
		return &Storage{pool: pool, timeout: timeout}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO sso.users (email, pass_hash) VALUES ($1, $2) RETURNING id`, email, passHash).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user models.User
	err := s.pool.QueryRow(ctx, 
		`SELECT id, email, pass_hash FROM sso.users WHERE email = $1`, email,
	).Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error)  {
	const op = "storage.posthres.IsAdmin"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var isAdmin bool
	err := s.pool.QueryRow(
		ctx,
		`SELECT is_admin from sso.users WHERE id = $1`,
		userId,
	).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err) 
		}
		return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var app models.App
	err := s.pool.QueryRow(
		ctx, `SELECT id, name, secret FROM sso.apps WHERE id = $1`,
		appID,
	).Scan(&app.Id, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, err
	}
	return app, nil
} 


func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"

	
}

