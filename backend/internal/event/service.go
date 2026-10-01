package event

import (
	"context"
	"strings"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Quote struct {
	ID        uuid.UUID  `json:"id"`
	EventType string     `json:"eventType"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Rooms     int        `json:"rooms"`
	People    int        `json:"people"`
	EventAt   *time.Time `json:"eventAt,omitempty"`
	Coffee    bool       `json:"coffee"`
	Lunch     bool       `json:"lunch"`
	Dinner    bool       `json:"dinner"`
	CreatedAt time.Time  `json:"createdAt"`
}

type CreateInput struct {
	EventType string `json:"eventType"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Rooms     int    `json:"rooms"`
	People    int    `json:"people"`
	EventAt   string `json:"eventAt"`
	Coffee    bool   `json:"coffee"`
	Lunch     bool   `json:"lunch"`
	Dinner    bool   `json:"dinner"`
}

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (Quote, error) {
	if strings.TrimSpace(in.EventType) == "" || strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Email) == "" {
		return Quote{}, apperr.New("INVALID_INPUT", "tipo, nome e e-mail são obrigatórios", apperr.ErrInvalidInput)
	}
	var eventAt *time.Time
	if strings.TrimSpace(in.EventAt) != "" {
		t, err := time.Parse(time.RFC3339, in.EventAt)
		if err != nil {
			t2, err2 := time.Parse("2006-01-02T15:04", in.EventAt)
			if err2 != nil {
				return Quote{}, apperr.New("INVALID_EVENT_AT", "data do evento inválida", apperr.ErrInvalidInput)
			}
			t = t2
		}
		eventAt = &t
	}
	var q Quote
	err := s.pool.QueryRow(ctx, `
		INSERT INTO event_quotes (event_type, name, email, phone, rooms, people, event_at, coffee, lunch, dinner)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, event_type, name, email, phone, rooms, people, event_at, coffee, lunch, dinner, created_at
	`, in.EventType, in.Name, in.Email, in.Phone, in.Rooms, in.People, eventAt, in.Coffee, in.Lunch, in.Dinner).Scan(
		&q.ID, &q.EventType, &q.Name, &q.Email, &q.Phone, &q.Rooms, &q.People, &q.EventAt, &q.Coffee, &q.Lunch, &q.Dinner, &q.CreatedAt,
	)
	return q, err
}
