package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Payload struct {
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(email string, username string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
func (p *Payload) Valid() error {
	// Check if the token is expired
	if time.Now().After(p.ExpiredAt) {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}

	// Add custom validation logic if needed

	return nil
}
