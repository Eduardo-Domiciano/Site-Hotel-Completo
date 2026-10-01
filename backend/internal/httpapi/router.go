package httpapi

import (
	"net/http"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/auth"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/career"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/event"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/reservation"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/room"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	Auth         *auth.Service
	Tokens       *auth.Tokens
	Rooms        *room.Repository
	Reservations *reservation.Service
	Events       *event.Service
	Careers      *career.Service
	CORSOrigins  []string
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORS(d.CORSOrigins))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		WriteData(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(api chi.Router) {
		authRL := RateLimit(time.Minute, 20)
		api.With(authRL).Post("/auth/register", registerHandler(d.Auth))
		api.With(authRL).Post("/auth/login", loginHandler(d.Auth))
		api.With(authRL).Post("/auth/refresh", refreshHandler(d.Auth))
		api.Post("/auth/logout", logoutHandler(d.Auth))

		api.Get("/rooms", listRoomsHandler(d.Rooms))
		api.Post("/event-quotes", createEventHandler(d.Events))
		api.Post("/job-applications", createJobHandler(d.Careers))
		api.Post("/reservations/guest", createGuestReservationHandler(d.Reservations))

		api.Group(func(priv chi.Router) {
			priv.Use(Auth(d.Tokens))
			priv.Get("/me", meHandler(d.Auth))
			priv.Put("/me", updateMeHandler(d.Auth))
			priv.Put("/me/password", changePasswordHandler(d.Auth))
			priv.Post("/reservations", createReservationHandler(d.Reservations))
			priv.Get("/reservations", listReservationsHandler(d.Reservations))
		})
	})

	return r
}
