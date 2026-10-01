package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/auth"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/career"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/config"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/db"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/event"
	httpapi "github.com/Eduardo-Domiciano/sitehotel-api/internal/httpapi"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/reservation"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/room"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/user"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := db.Migrate(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		slog.Error("migrate", "err", err)
		os.Exit(1)
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		slog.Error("upload dir", "err", err)
		os.Exit(1)
	}

	userRepo := user.NewRepository(pool)
	tokens := auth.NewTokens(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	authSvc := auth.NewService(userRepo, tokens)
	roomRepo := room.NewRepository(pool)
	reservationSvc := reservation.NewService(pool, roomRepo)
	eventSvc := event.NewService(pool)
	careerSvc := career.NewService(pool, cfg.UploadDir, cfg.MaxUploadBytes)

	router := httpapi.NewRouter(httpapi.Deps{
		Auth:         authSvc,
		Tokens:       tokens,
		Rooms:        roomRepo,
		Reservations: reservationSvc,
		Events:       eventSvc,
		Careers:      careerSvc,
		CORSOrigins:  cfg.CORSOrigins,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("api listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
