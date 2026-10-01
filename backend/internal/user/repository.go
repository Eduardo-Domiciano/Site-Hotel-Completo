package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, name, email, passwordHash string, phone *string, birthDate *time.Time) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, phone, birth_date)
		VALUES ($1, LOWER($2), $3, $4, $5)
		RETURNING id, name, email, password_hash, phone, birth_date, avatar_url, created_at, updated_at
	`, name, email, passwordHash, phone, birthDate).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.BirthDate, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, apperr.New("EMAIL_TAKEN", "e-mail já cadastrado", apperr.ErrConflict)
		}
		return User{}, err
	}
	return u, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, phone, birth_date, avatar_url, created_at, updated_at
		FROM users WHERE LOWER(email) = LOWER($1)
	`, strings.TrimSpace(email)).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.BirthDate, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	return u, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, phone, birth_date, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.BirthDate, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	return u, err
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, name string, phone *string, birthDate *time.Time) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `
		UPDATE users
		SET name = $2, phone = $3, birth_date = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, password_hash, phone, birth_date, avatar_url, created_at, updated_at
	`, id, name, phone, birthDate).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Phone, &u.BirthDate, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	return u, err
}

func (r *Repository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1
	`, id, passwordHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	return nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *Repository) FindRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, *time.Time, error) {
	var userID uuid.UUID
	var expiresAt time.Time
	var revokedAt *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = $1
	`, tokenHash).Scan(&userID, &expiresAt, &revokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, time.Time{}, nil, apperr.New("TOKEN_INVALID", "refresh token inválido", apperr.ErrUnauthorized)
	}
	return userID, expiresAt, revokedAt, err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL
	`, tokenHash)
	return err
}

func (r *Repository) RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	return err
}
