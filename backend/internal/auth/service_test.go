package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/auth"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/user"
	"github.com/google/uuid"
)

type memStore struct {
	byEmail map[string]user.User
	byID    map[uuid.UUID]user.User
}

func newMem() *memStore {
	return &memStore{byEmail: map[string]user.User{}, byID: map[uuid.UUID]user.User{}}
}

func (m *memStore) Create(_ context.Context, name, email, passwordHash string, phone *string, birthDate *time.Time) (user.User, error) {
	if _, ok := m.byEmail[email]; ok {
		return user.User{}, apperr.New("EMAIL_TAKEN", "e-mail já cadastrado", apperr.ErrConflict)
	}
	u := user.User{
		ID: uuid.New(), Name: name, Email: email, PasswordHash: passwordHash,
		Phone: phone, BirthDate: birthDate, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	m.byEmail[email] = u
	m.byID[u.ID] = u
	return u, nil
}
func (m *memStore) FindByEmail(_ context.Context, email string) (user.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return user.User{}, apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	return u, nil
}
func (m *memStore) FindByID(_ context.Context, id uuid.UUID) (user.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return user.User{}, apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	return u, nil
}
func (m *memStore) UpdateProfile(_ context.Context, id uuid.UUID, name string, phone *string, birthDate *time.Time) (user.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return user.User{}, apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	u.Name, u.Phone, u.BirthDate = name, phone, birthDate
	m.byID[id] = u
	m.byEmail[u.Email] = u
	return u, nil
}
func (m *memStore) UpdatePassword(_ context.Context, id uuid.UUID, passwordHash string) error {
	u, ok := m.byID[id]
	if !ok {
		return apperr.New("USER_NOT_FOUND", "usuário não encontrado", apperr.ErrNotFound)
	}
	u.PasswordHash = passwordHash
	m.byID[id] = u
	m.byEmail[u.Email] = u
	return nil
}
func (m *memStore) SaveRefreshToken(context.Context, uuid.UUID, string, time.Time) error { return nil }
func (m *memStore) FindRefreshToken(context.Context, string) (uuid.UUID, time.Time, *time.Time, error) {
	return uuid.Nil, time.Time{}, nil, apperr.New("TOKEN_INVALID", "refresh token inválido", apperr.ErrUnauthorized)
}
func (m *memStore) RevokeRefreshToken(context.Context, string) error             { return nil }
func (m *memStore) RevokeAllRefreshTokens(context.Context, uuid.UUID) error       { return nil }

func TestRegisterAndLogin(t *testing.T) {
	svc := auth.NewService(newMem(), auth.NewTokens("secret", time.Minute, time.Hour))
	reg, err := svc.Register(context.Background(), auth.RegisterInput{
		Name: "Eduardo", Email: "edu@hotel.com", Password: "senha1234",
	})
	if err != nil {
		t.Fatal(err)
	}
	if reg.User.Email != "edu@hotel.com" || reg.Tokens.AccessToken == "" {
		t.Fatalf("bad register result: %+v", reg)
	}
	login, err := svc.Login(context.Background(), auth.LoginInput{
		Email: "edu@hotel.com", Password: "senha1234",
	})
	if err != nil {
		t.Fatal(err)
	}
	if login.User.ID != reg.User.ID {
		t.Fatal("login user mismatch")
	}
}
