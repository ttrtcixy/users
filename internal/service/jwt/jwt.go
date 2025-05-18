package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"strconv"
	"time"
)

type Hasher interface {
	Hash(str string) (hash string, err error)
}

type JwtTokenService struct {
	cfg *config.JWTConfig
}

func New(cfg *config.JWTConfig) *JwtTokenService {
	return &JwtTokenService{
		cfg,
	}
}

type EmailVerificationClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleId   string `json:"role_id"`
	jwt.RegisteredClaims
}

// ParseVerificationToken - UserError: apperrors.ErrEmailTokenExpired
func (t *JwtTokenService) ParseVerificationToken(jwtToken string) (email string, err error) {
	const op = "JwtTokenService.Parse"
	// parse and validate token
	token, err := jwt.ParseWithClaims(jwtToken, &EmailVerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.cfg.JWTSecret()), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", apperrors.ErrEmailTokenExpired
		}
		return "", err
	}

	// get user email to activate account
	claims, ok := token.Claims.(*EmailVerificationClaims)
	if !ok {
		return "", fmt.Errorf("%s: token structure invalid or signature incorrect", op)
	}

	if claims.Email == "" {
		return "", fmt.Errorf("%s: email cannot be empty", op)
	}

	return claims.Email, nil
}

// AccessToken - generate access token with user username, email, roleId.
func (t *JwtTokenService) AccessToken(user *entities.User) (token string, err error) {
	const op = "JwtTokenService.AccessToken"

	exp := time.Now().Add(t.cfg.AccessJwtExpiry())

	claims := &UserClaims{
		Username: user.Username,
		Email:    user.Email,
		RoleId:   strconv.Itoa(user.RoleId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "auth_grpc_app",
		},
	}

	if token, err = t.jwt(t.cfg.JWTSecret(), claims); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// todo test refresh token parse

// RefreshToken - generate jwtToken and hash it to create user session.
func (t *JwtTokenService) RefreshToken() (token string, err error) {
	const op = "JwtTokenService.RefreshToken"

	exp := time.Now().Add(t.cfg.RefreshJwtExpiry())

	claims := &jwt.RegisteredClaims{
		Issuer:    "auth_grpc_app",
		ExpiresAt: jwt.NewNumericDate(exp),
	}

	if token, err = t.jwt(t.cfg.JWTSecret(), claims); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// EmailVerificationToken - generate jwtToken with user email to verify.
func (t *JwtTokenService) EmailVerificationToken(email string) (token string, err error) {
	const op = "JwtTokenService.EmailVerificationToken"

	expAt := time.Now().Add(t.cfg.EmailJwtExpiry())

	// todo test pointer
	claims := &EmailVerificationClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth_grpc_app",
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	}

	if token, err = t.jwt(t.cfg.JWTSecret(), claims); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (t *JwtTokenService) jwt(secret string, claims jwt.Claims) (token string, err error) {
	const op = "JwtTokenService.jwt"
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if token, err = jwtToken.SignedString([]byte(secret)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
