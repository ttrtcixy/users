package ports

import (
	"github.com/ttrtcixy/users/internal/core/entities"
	"time"
)

type JwtService interface {
	ParseVerificationToken(jwtToken string) (string, error)
	EmailVerificationToken(email string) (token string, err error)

	AccessToken(user *entities.User) (token string, err error)
	RefreshToken(clientID, tokenID string, exp time.Time) (token string, err error)
}
