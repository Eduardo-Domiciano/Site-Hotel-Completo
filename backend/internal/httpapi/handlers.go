package httpapi

import (
	"net/http"
	"strings"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/auth"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/career"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/event"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/reservation"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/room"
)

func registerHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in auth.RegisterInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.Register(r.Context(), in)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusCreated, out)
	}
}

func loginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in auth.LoginInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.Login(r.Context(), in)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusOK, out)
	}
}

func refreshHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			RefreshToken string `json:"refreshToken"`
		}
		if err := DecodeJSON(r, &body); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.Refresh(r.Context(), body.RefreshToken)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusOK, out)
	}
}

func logoutHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			RefreshToken string `json:"refreshToken"`
		}
		_ = DecodeJSON(r, &body)
		_ = svc.Logout(r.Context(), body.RefreshToken)
		WriteData(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func meHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		if !ok {
			WriteError(w, apperr.New("UNAUTHORIZED", "não autenticado", apperr.ErrUnauthorized))
			return
		}
		out, err := svc.Me(r.Context(), id)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusOK, out)
	}
}

func updateMeHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		if !ok {
			WriteError(w, apperr.New("UNAUTHORIZED", "não autenticado", apperr.ErrUnauthorized))
			return
		}
		var in auth.UpdateProfileInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.UpdateProfile(r.Context(), id, in)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusOK, out)
	}
}

func changePasswordHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		if !ok {
			WriteError(w, apperr.New("UNAUTHORIZED", "não autenticado", apperr.ErrUnauthorized))
			return
		}
		var in auth.ChangePasswordInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		if err := svc.ChangePassword(r.Context(), id, in); err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func listRoomsHandler(repo *room.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := repo.List(r.Context())
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusOK, rooms)
	}
}

func createReservationHandler(svc *reservation.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		if !ok {
			WriteError(w, apperr.New("UNAUTHORIZED", "não autenticado", apperr.ErrUnauthorized))
			return
		}
		var in reservation.CreateInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		uid := id
		out, err := svc.Create(r.Context(), &uid, in)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusCreated, out)
	}
}

func createGuestReservationHandler(svc *reservation.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in reservation.CreateInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.Create(r.Context(), nil, in)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusCreated, out)
	}
}

func listReservationsHandler(svc *reservation.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		if !ok {
			WriteError(w, apperr.New("UNAUTHORIZED", "não autenticado", apperr.ErrUnauthorized))
			return
		}
		out, err := svc.ListByUser(r.Context(), id)
		if err != nil {
			WriteError(w, err)
			return
		}
		if out == nil {
			out = []reservation.Reservation{}
		}
		WriteData(w, http.StatusOK, out)
	}
}

func createEventHandler(svc *event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in event.CreateInput
		if err := DecodeJSON(r, &in); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.Create(r.Context(), in)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusCreated, out)
	}
}

func createJobHandler(svc *career.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				WriteError(w, apperr.New("INVALID_FORM", "formulário inválido", apperr.ErrInvalidInput))
				return
			}
			area := r.FormValue("area")
			name := r.FormValue("name")
			email := r.FormValue("email")
			phone := r.FormValue("phone")
			var phonePtr *string
			if strings.TrimSpace(phone) != "" {
				phonePtr = &phone
			}
			file, header, err := r.FormFile("resume")
			var filename string
			var size int64
			if err == nil {
				defer file.Close()
				filename = header.Filename
				size = header.Size
				out, err := svc.Create(r.Context(), area, name, email, phonePtr, filename, file, size)
				if err != nil {
					WriteError(w, err)
					return
				}
				WriteData(w, http.StatusCreated, out)
				return
			}
			out, err := svc.Create(r.Context(), area, name, email, phonePtr, "", nil, 0)
			if err != nil {
				WriteError(w, err)
				return
			}
			WriteData(w, http.StatusCreated, out)
			return
		}

		var body struct {
			Area  string  `json:"area"`
			Name  string  `json:"name"`
			Email string  `json:"email"`
			Phone *string `json:"phone"`
		}
		if err := DecodeJSON(r, &body); err != nil {
			WriteError(w, err)
			return
		}
		out, err := svc.Create(r.Context(), body.Area, body.Name, body.Email, body.Phone, "", nil, 0)
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteData(w, http.StatusCreated, out)
	}
}
