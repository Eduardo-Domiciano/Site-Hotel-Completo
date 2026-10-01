package httpapi

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/auth"
	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

func CORS(origins []string) func(http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, o := range origins {
		allowed[o] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Auth(tokens *auth.Tokens) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				WriteError(w, apperr.New("UNAUTHORIZED", "token ausente", apperr.ErrUnauthorized))
				return
			}
			claims, err := tokens.ParseAccess(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				WriteError(w, apperr.New("UNAUTHORIZED", "token inválido", apperr.ErrUnauthorized))
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type rateLimiter struct {
	mu      sync.Mutex
	hits    map[string][]time.Time
	window  time.Duration
	maxHits int
}

func newRateLimiter(window time.Duration, maxHits int) *rateLimiter {
	return &rateLimiter{hits: map[string][]time.Time{}, window: window, maxHits: maxHits}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-rl.window)
	list := rl.hits[key]
	kept := list[:0]
	for _, t := range list {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= rl.maxHits {
		rl.hits[key] = kept
		return false
	}
	rl.hits[key] = append(kept, now)
	return true
}

func RateLimit(window time.Duration, maxHits int) func(http.Handler) http.Handler {
	rl := newRateLimiter(window, maxHits)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
				key = strings.Split(xf, ",")[0]
			}
			key += "|" + r.URL.Path
			if !rl.allow(key) {
				WriteError(w, apperr.New("RATE_LIMITED", "muitas tentativas, aguarde", apperr.ErrTooManyRequests))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
