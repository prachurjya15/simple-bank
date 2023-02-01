package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("Token is expired and cant be used")
	ErrInvalidToken = errors.New("Invalid token provided")
)

type Payload struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	IssuedAt       time.Time `json:"issued_at"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	p := Payload{
		ID:             tokenId,
		Username:       username,
		IssuedAt:       time.Now(),
		ExpirationTime: time.Now().Add(duration),
	}

	return &p, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpirationTime) {
		return ErrExpiredToken
	}
	return nil
}
