package entities

import "time"

type CreateSession struct {
	UserID           int
	ClientUUID       string
	RefreshTokenUUID string
	ExpiresAt        time.Time
}

type RefreshSession struct {
	ClientUUID          string
	OldRefreshTokenUUID string
	NewRefreshTokenUUID string
	ExpiresAt           time.Time
}

type TokenUserInfo struct {
	ID       int
	Username string
	Email    string
	RoleID   int
}
