package reservation

import (
	"context"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/Eduardo-Domiciano/sitehotel-api/internal/room"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type RoomLine struct {
	RoomTypeID   string          `json:"roomTypeId"`
	Adults       int             `json:"adults"`
	Children     int             `json:"children"`
	NightlyPrice decimal.Decimal `json:"nightlyPrice"`
	LineTotal    decimal.Decimal `json:"lineTotal"`
}

type Reservation struct {
	ID         uuid.UUID       `json:"id"`
	UserID     *uuid.UUID      `json:"userId,omitempty"`
	GuestName  string          `json:"guestName"`
	GuestPhone string          `json:"guestPhone"`
	GuestEmail string          `json:"guestEmail"`
	CheckIn    string          `json:"checkIn"`
	CheckOut   string          `json:"checkOut"`
	Nights     int             `json:"nights"`
	Total      decimal.Decimal `json:"total"`
	Status     string          `json:"status"`
	Rooms      []RoomLine      `json:"rooms"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type CreateRoomInput struct {
	RoomTypeID string `json:"roomTypeId"`
	Adults     int    `json:"adults"`
	Children   int    `json:"children"`
}

type CreateInput struct {
	GuestName  string            `json:"guestName"`
	GuestPhone string            `json:"guestPhone"`
	GuestEmail string            `json:"guestEmail"`
	CheckIn    string            `json:"checkIn"`
	CheckOut   string            `json:"checkOut"`
	Rooms      []CreateRoomInput `json:"rooms"`
}

type Service struct {
	pool  *pgxpool.Pool
	rooms *room.Repository
}

func NewService(pool *pgxpool.Pool, rooms *room.Repository) *Service {
	return &Service{pool: pool, rooms: rooms}
}

func (s *Service) Create(ctx context.Context, userID *uuid.UUID, in CreateInput) (Reservation, error) {
	if len(in.Rooms) == 0 || len(in.Rooms) > 3 {
		return Reservation{}, apperr.New("INVALID_ROOMS", "selecione entre 1 e 3 tipos de quarto", apperr.ErrInvalidInput)
	}
	checkIn, err := time.Parse("2006-01-02", in.CheckIn)
	if err != nil {
		return Reservation{}, apperr.New("INVALID_CHECK_IN", "check-in inválido", apperr.ErrInvalidInput)
	}
	checkOut, err := time.Parse("2006-01-02", in.CheckOut)
	if err != nil {
		return Reservation{}, apperr.New("INVALID_CHECK_OUT", "check-out inválido", apperr.ErrInvalidInput)
	}
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	if nights < 1 {
		return Reservation{}, apperr.New("INVALID_DATES", "check-out deve ser após o check-in", apperr.ErrInvalidInput)
	}
	if in.GuestName == "" || in.GuestPhone == "" || in.GuestEmail == "" {
		return Reservation{}, apperr.New("INVALID_GUEST", "preencha nome, telefone e e-mail", apperr.ErrInvalidInput)
	}

	seen := map[string]bool{}
	var lines []RoomLine
	total := decimal.Zero
	for _, item := range in.Rooms {
		if seen[item.RoomTypeID] {
			return Reservation{}, apperr.New("DUPLICATE_ROOM", "tipos de quarto duplicados", apperr.ErrInvalidInput)
		}
		seen[item.RoomTypeID] = true
		rt, err := s.rooms.Get(ctx, item.RoomTypeID)
		if err != nil {
			return Reservation{}, apperr.New("ROOM_NOT_FOUND", "quarto não encontrado", apperr.ErrNotFound)
		}
		if item.Adults < 1 || item.Adults > rt.MaxAdults {
			return Reservation{}, apperr.New("INVALID_ADULTS", "número de adultos inválido", apperr.ErrInvalidInput)
		}
		if item.Children < 0 || item.Children > rt.MaxChildren {
			return Reservation{}, apperr.New("INVALID_CHILDREN", "número de crianças inválido", apperr.ErrInvalidInput)
		}
		occ := item.Adults
		if occ > 4 {
			occ = 4
		}
		price, ok := rt.Prices[itoa(occ)]
		if !ok {
			price, ok = rt.Prices["1"]
			if !ok {
				return Reservation{}, apperr.New("PRICE_MISSING", "preço indisponível", apperr.ErrInvalidInput)
			}
		}
		lineTotal := price.To.Mul(decimal.NewFromInt(int64(nights)))
		lines = append(lines, RoomLine{
			RoomTypeID:   rt.ID,
			Adults:       item.Adults,
			Children:     item.Children,
			NightlyPrice: price.To,
			LineTotal:    lineTotal,
		})
		total = total.Add(lineTotal)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Reservation{}, err
	}
	defer tx.Rollback(ctx)

	var res Reservation
	err = tx.QueryRow(ctx, `
		INSERT INTO reservations (user_id, guest_name, guest_phone, guest_email, check_in, check_out, nights, total)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, user_id, guest_name, guest_phone, guest_email, check_in::text, check_out::text, nights, total, status, created_at
	`, userID, in.GuestName, in.GuestPhone, in.GuestEmail, checkIn, checkOut, nights, total).Scan(
		&res.ID, &res.UserID, &res.GuestName, &res.GuestPhone, &res.GuestEmail, &res.CheckIn, &res.CheckOut, &res.Nights, &res.Total, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return Reservation{}, err
	}

	for _, line := range lines {
		_, err := tx.Exec(ctx, `
			INSERT INTO reservation_rooms (reservation_id, room_type_id, adults, children, nightly_price, line_total)
			VALUES ($1,$2,$3,$4,$5,$6)
		`, res.ID, line.RoomTypeID, line.Adults, line.Children, line.NightlyPrice, line.LineTotal)
		if err != nil {
			return Reservation{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return Reservation{}, err
	}
	res.Rooms = lines
	return res, nil
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]Reservation, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, guest_name, guest_phone, guest_email, check_in::text, check_out::text, nights, total, status, created_at
		FROM reservations WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Reservation
	for rows.Next() {
		var r Reservation
		if err := rows.Scan(&r.ID, &r.UserID, &r.GuestName, &r.GuestPhone, &r.GuestEmail, &r.CheckIn, &r.CheckOut, &r.Nights, &r.Total, &r.Status, &r.CreatedAt); err != nil {
			return nil, err
		}
		rooms, err := s.loadRooms(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Rooms = rooms
		list = append(list, r)
	}
	return list, rows.Err()
}

func (s *Service) loadRooms(ctx context.Context, reservationID uuid.UUID) ([]RoomLine, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT room_type_id, adults, children, nightly_price, line_total
		FROM reservation_rooms WHERE reservation_id = $1
	`, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []RoomLine
	for rows.Next() {
		var line RoomLine
		if err := rows.Scan(&line.RoomTypeID, &line.Adults, &line.Children, &line.NightlyPrice, &line.LineTotal); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, rows.Err()
}

func itoa(n int) string {
	switch n {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	default:
		return "4"
	}
}
