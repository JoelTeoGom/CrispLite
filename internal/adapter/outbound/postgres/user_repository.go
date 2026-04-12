package postgres

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool   *pgxpool.Pool
	logger outbound.Logger
}

func NewUserRepo(pool *pgxpool.Pool, logger outbound.Logger) *UserRepo {
	return &UserRepo{pool: pool, logger: logger}
}

func (p *UserRepo) Save(ctx context.Context, user *domain.User) (string, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	userID := uuid.New().String()
	_, err = tx.Exec(ctx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", userID, user.Username, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", domain.ErrUserAlreadyExists
		}
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return userID, nil
}

func (p *UserRepo) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := p.pool.QueryRow(ctx,
		"SELECT id, username, password, created_at FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return &user, nil
}

func (p *UserRepo) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := p.pool.QueryRow(ctx,
		"SELECT id, username, password, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return &user, nil
}

func (p *UserRepo) AddContact(ctx context.Context, userID, contactID string) error {
	_, err := p.pool.Exec(ctx, "INSERT INTO contacts (user_id, contact_id) VALUES ($1, $2)", userID, contactID)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		switch pgErr.Code {
		case "23503":
			return domain.ErrUserNotFound
		case "23505":
			return domain.ErrContactAlreadyExists
		default:
			return err
		}
	}
	return nil
}

func (p *UserRepo) SearchUsers(ctx context.Context, query string, limit, offset int) ([]domain.UserSummary, error) {
	rows, err := p.pool.Query(ctx,
		"SELECT id, username FROM users WHERE username ILIKE $1 ORDER BY username LIMIT $2 OFFSET $3",
		"%"+query+"%", limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserSummary
	for rows.Next() {
		var u domain.UserSummary
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if users == nil {
		users = []domain.UserSummary{}
	}
	return users, nil
}

func (p *UserRepo) RemoveContact(ctx context.Context, userID, contactID string) error {
	_, err := p.pool.Exec(ctx, "DELETE FROM contacts WHERE user_id = $1 AND contact_id = $2", userID, contactID)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23503" {
			return domain.ErrUserNotFound
		}
	}
	return err
}
