package career

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	ID         uuid.UUID `json:"id"`
	Area       string    `json:"area"`
	Name       string    `json:"name"`
	Phone      *string   `json:"phone,omitempty"`
	Email      string    `json:"email"`
	ResumePath *string   `json:"resumePath,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Service struct {
	pool      *pgxpool.Pool
	uploadDir string
	maxBytes  int64
}

func NewService(pool *pgxpool.Pool, uploadDir string, maxBytes int64) *Service {
	return &Service{pool: pool, uploadDir: uploadDir, maxBytes: maxBytes}
}

func (s *Service) Create(ctx context.Context, area, name, email string, phone *string, filename string, src io.Reader, size int64) (Application, error) {
	if strings.TrimSpace(area) == "" || strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
		return Application{}, apperr.New("INVALID_INPUT", "área, nome e e-mail são obrigatórios", apperr.ErrInvalidInput)
	}
	var resumePath *string
	if src != nil && filename != "" {
		if size > s.maxBytes {
			return Application{}, apperr.New("FILE_TOO_LARGE", "currículo excede o tamanho máximo", apperr.ErrInvalidInput)
		}
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".pdf", ".doc", ".docx":
		default:
			return Application{}, apperr.New("INVALID_FILE", "envie PDF, DOC ou DOCX", apperr.ErrInvalidInput)
		}
		if err := os.MkdirAll(filepath.Join(s.uploadDir, "resumes"), 0o755); err != nil {
			return Application{}, err
		}
		id := uuid.NewString()
		rel := filepath.Join("resumes", id+ext)
		dstPath := filepath.Join(s.uploadDir, rel)
		f, err := os.Create(dstPath)
		if err != nil {
			return Application{}, err
		}
		defer f.Close()
		if _, err := io.Copy(f, io.LimitReader(src, s.maxBytes+1)); err != nil {
			return Application{}, err
		}
		resumePath = &rel
	}

	var app Application
	err := s.pool.QueryRow(ctx, `
		INSERT INTO job_applications (area, name, phone, email, resume_path)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, area, name, phone, email, resume_path, created_at
	`, area, name, phone, email, resumePath).Scan(
		&app.ID, &app.Area, &app.Name, &app.Phone, &app.Email, &app.ResumePath, &app.CreatedAt,
	)
	return app, err
}
