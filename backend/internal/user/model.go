package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Phone        *string    `json:"phone,omitempty"`
	BirthDate    *time.Time `json:"birthDate,omitempty"`
	AvatarURL    *string    `json:"avatarUrl,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type PublicUser struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone,omitempty"`
	BirthDate *string    `json:"birthDate,omitempty"`
	AvatarURL *string    `json:"avatarUrl,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u User) Public() PublicUser {
	var birth *string
	if u.BirthDate != nil {
		s := u.BirthDate.Format("2006-01-02")
		birth = &s
	}
	return PublicUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		BirthDate: birth,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
