package auth

import (
	"context"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/user"
	"github.com/google/uuid"
)

type UserStore interface {
	Create(ctx context.Context, name, email, passwordHash string, phone *string, birthDate *time.Time) (user.User, error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (user.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, name string, phone *string, birthDate *time.Time) (user.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, *time.Time, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type Service struct {
	users  UserStore
	tokens *Tokens
}

func NewService(users UserStore, tokens *Tokens) *Service {
	return &Service{users: users, tokens: tokens}
}

type RegisterInput struct {
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Phone     *string `json:"phone"`
	BirthDate *string `json:"birthDate"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileInput struct {
	Name      string  `json:"name"`
	Phone     *string `json:"phone"`
	BirthDate *string `json:"birthDate"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type AuthResult struct {
	User   user.PublicUser `json:"user"`
	Tokens TokenPair       `json:"tokens"`
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (AuthResult, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if name == "" || utf8.RuneCountInString(name) < 2 {
		return AuthResult{}, apperr.New("INVALID_NAME", "nome inválido", apperr.ErrInvalidInput)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return AuthResult{}, apperr.New("INVALID_EMAIL", "e-mail inválido", apperr.ErrInvalidInput)
	}
	if len(in.Password) < 8 {
		return AuthResult{}, apperr.New("INVALID_PASSWORD", "senha deve ter pelo menos 8 caracteres", apperr.ErrInvalidInput)
	}
	birth, err := parseDatePtr(in.BirthDate)
	if err != nil {
		return AuthResult{}, apperr.New("INVALID_BIRTH_DATE", "data de nascimento inválida", apperr.ErrInvalidInput)
	}
	hash, err := HashPassword(in.Password)
	if err != nil {
		return AuthResult{}, err
	}
	u, err := s.users.Create(ctx, name, email, hash, emptyToNil(in.Phone), birth)
	if err != nil {
		return AuthResult{}, err
	}
	return s.issueFor(ctx, u)
}

func (s *Service) Login(ctx context.Context, in LoginInput) (AuthResult, error) {
	u, err := s.users.FindByEmail(ctx, in.Email)
	if err != nil {
		return AuthResult{}, apperr.New("INVALID_CREDENTIALS", "e-mail ou senha inválidos", apperr.ErrUnauthorized)
	}
	if !CheckPassword(u.PasswordHash, in.Password) {
		return AuthResult{}, apperr.New("INVALID_CREDENTIALS", "e-mail ou senha inválidos", apperr.ErrUnauthorized)
	}
	return s.issueFor(ctx, u)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	hash := HashToken(refreshToken)
	userID, expiresAt, revokedAt, err := s.users.FindRefreshToken(ctx, hash)
	if err != nil {
		return AuthResult{}, err
	}
	if revokedAt != nil || time.Now().UTC().After(expiresAt) {
		return AuthResult{}, apperr.New("TOKEN_EXPIRED", "refresh token expirado", apperr.ErrUnauthorized)
	}
	_ = s.users.RevokeRefreshToken(ctx, hash)
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return AuthResult{}, err
	}
	return s.issueFor(ctx, u)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.users.RevokeRefreshToken(ctx, HashToken(refreshToken))
}

func (s *Service) Me(ctx context.Context, userID uuid.UUID) (user.PublicUser, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return user.PublicUser{}, err
	}
	return u.Public(), nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, in UpdateProfileInput) (user.PublicUser, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) < 2 {
		return user.PublicUser{}, apperr.New("INVALID_NAME", "nome inválido", apperr.ErrInvalidInput)
	}
	birth, err := parseDatePtr(in.BirthDate)
	if err != nil {
		return user.PublicUser{}, apperr.New("INVALID_BIRTH_DATE", "data de nascimento inválida", apperr.ErrInvalidInput)
	}
	u, err := s.users.UpdateProfile(ctx, userID, name, emptyToNil(in.Phone), birth)
	if err != nil {
		return user.PublicUser{}, err
	}
	return u.Public(), nil
}

func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, in ChangePasswordInput) error {
	if len(in.NewPassword) < 8 {
		return apperr.New("INVALID_PASSWORD", "senha deve ter pelo menos 8 caracteres", apperr.ErrInvalidInput)
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !CheckPassword(u.PasswordHash, in.CurrentPassword) {
		return apperr.New("INVALID_CREDENTIALS", "senha atual incorreta", apperr.ErrUnauthorized)
	}
	hash, err := HashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}
	return s.users.RevokeAllRefreshTokens(ctx, userID)
}

func (s *Service) issueFor(ctx context.Context, u user.User) (AuthResult, error) {
	pair, refreshHash, refreshExp, err := s.tokens.Issue(u.ID, u.Email)
	if err != nil {
		return AuthResult{}, err
	}
	if err := s.users.SaveRefreshToken(ctx, u.ID, refreshHash, refreshExp); err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: u.Public(), Tokens: pair}, nil
}

func parseDatePtr(v *string) (*time.Time, error) {
	if v == nil || strings.TrimSpace(*v) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*v))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func emptyToNil(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}
