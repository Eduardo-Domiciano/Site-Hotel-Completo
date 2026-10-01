package room

import (
	"context"
	"strconv"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/apperr"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PriceTier map[int]Price

type Price struct {
	From decimal.Decimal `json:"from"`
	To   decimal.Decimal `json:"to"`
}

type Type struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	MaxAdults   int            `json:"maxAdults"`
	MaxChildren int            `json:"maxChildren"`
	Prices      map[string]Price `json:"prices"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context) ([]Type, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, image, max_adults, max_children
		FROM room_types
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Type
	index := map[string]int{}
	for rows.Next() {
		var t Type
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Image, &t.MaxAdults, &t.MaxChildren); err != nil {
			return nil, err
		}
		t.Prices = map[string]Price{}
		index[t.ID] = len(rooms)
		rooms = append(rooms, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	priceRows, err := r.pool.Query(ctx, `
		SELECT room_type_id, occupancy, price_from, price_to FROM room_prices
	`)
	if err != nil {
		return nil, err
	}
	defer priceRows.Close()
	for priceRows.Next() {
		var id string
		var occ int
		var from, to decimal.Decimal
		if err := priceRows.Scan(&id, &occ, &from, &to); err != nil {
			return nil, err
		}
		i, ok := index[id]
		if !ok {
			continue
		}
		rooms[i].Prices[itoa(occ)] = Price{From: from, To: to}
	}
	return rooms, priceRows.Err()
}

func (r *Repository) Get(ctx context.Context, id string) (Type, error) {
	rooms, err := r.List(ctx)
	if err != nil {
		return Type{}, err
	}
	for _, room := range rooms {
		if room.ID == id {
			return room, nil
		}
	}
	return Type{}, apperr.New("ROOM_NOT_FOUND", "quarto não encontrado", apperr.ErrNotFound)
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
