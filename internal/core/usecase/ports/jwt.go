package ports

import "github.com/ttrtcixy/users/internal/core/entities"

type JwtService interface {
	ParseVerificationToken(jwtToken string) (string, error)
	AccessToken(user *entities.User) (token string, err error)
	RefreshToken() (jwtRefresh string, err error)
	EmailVerificationToken(email string) (token string, err error)
}
