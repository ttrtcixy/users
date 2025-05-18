package entities

import "time"

type CreateSession struct {
	UserID           int64
	RefreshTokenHash string
	ClientUUID       string
	RefreshTokenUUID string
	ExpiresAt        time.Time
}
