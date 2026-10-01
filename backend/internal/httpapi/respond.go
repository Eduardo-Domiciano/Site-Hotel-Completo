package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
)

type envelope struct {
	Data  any `json:"data,omitempty"`
	Error any `json:"error,omitempty"`
}

type errBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteData(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Data: data})
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "INTERNAL"
	message := "erro interno"

	var ae *apperr.Error
	if errors.As(err, &ae) {
		code = ae.Code
		message = ae.Message
		switch {
		case errors.Is(ae.Err, apperr.ErrInvalidInput):
			status = http.StatusBadRequest
		case errors.Is(ae.Err, apperr.ErrUnauthorized):
			status = http.StatusUnauthorized
		case errors.Is(ae.Err, apperr.ErrForbidden):
			status = http.StatusForbidden
		case errors.Is(ae.Err, apperr.ErrNotFound):
			status = http.StatusNotFound
		case errors.Is(ae.Err, apperr.ErrConflict):
			status = http.StatusConflict
		case errors.Is(ae.Err, apperr.ErrTooManyRequests):
			status = http.StatusTooManyRequests
		}
	} else {
		slog.Error("unhandled error", "err", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Error: errBody{Code: code, Message: message}})
}

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return apperr.New("INVALID_JSON", "JSON inválido", apperr.ErrInvalidInput)
	}
	return nil
}
