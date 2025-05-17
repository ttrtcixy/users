package authusecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"strconv"
	"time"
)

func (u *VerifyUseCase) jwt(secret string, claims jwt.Claims) (token string, err error) {
	const op = "JWT"
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if token, err = jwtToken.SignedString([]byte(secret)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

type VerifyUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo usecaseports.VerifyRepository
}

func NewVerify(ctx context.Context, cfg *config.UsecaseConfig, log logger.Logger, repo usecaseports.Repository) *VerifyUseCase {
	return &VerifyUseCase{
		cfg:  cfg,
		log:  log,
		repo: repo,
	}
}

type UserClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleId   string `json:"role_id"`
	jwt.RegisteredClaims
}

// todo нужна ли тут транзакция?

// Verify - get jwtToken with email and activate user with that email.
func (u *VerifyUseCase) Verify(ctx context.Context, payload *entities.VerifyRequest) (result *entities.VerifyResponse, err error) {
	const op = "VerifyUseCase.Verify"
	defer func() {
		if err != nil {
			if errors.Is(err, apperrors.ErrEmailTokenExpired) {
				return
			}
			u.log.ErrorOp(op, err)
			err = apperrors.ErrServer
		}
	}()

	email, err := u.parseJwt(payload.JwtToken)
	if err != nil {
		return nil, err
	}

	user, err := u.repo.ActivateUser(ctx, email)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.accessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenHash, err := u.refreshToken()
	if err != nil {
		return nil, err
	}

	if err = u.repo.CreateSession(ctx, user.ID, refreshTokenHash); err != nil {
		return nil, err
	}

	result = &entities.VerifyResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}

// parseJwt - UserError: apperrors.ErrEmailTokenExpired
func (u *VerifyUseCase) parseJwt(jwtToken string) (email string, err error) {
	const op = "parseJwt"
	// parse and validate token
	token, err := jwt.ParseWithClaims(jwtToken, &EmailVerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.cfg.JWTSecret()), nil
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

	t, _ := claims.GetExpirationTime()
	u.log.Debug(t.String())
	return claims.Email, nil
}

func (u *VerifyUseCase) accessToken(user *entities.User) (token string, err error) {
	const op = "accessToken"

	exp := time.Now().Add(u.cfg.AccessJwtExpiry())

	claims := &UserClaims{
		Username: user.Username,
		Email:    user.Email,
		RoleId:   strconv.Itoa(user.RoleId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "auth_grpc_app",
		},
	}

	if token, err = u.jwt(u.cfg.JWTSecret(), claims); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// todo test refresh token parse

// refreshToken - generate jwtToken and hash it to create session user session.
func (u *VerifyUseCase) refreshToken() (jwtRefresh, hash string, err error) {
	const op = "refreshToken"

	exp := time.Now().Add(u.cfg.RefreshJwtExpiry())

	claims := &jwt.RegisteredClaims{
		Issuer:    "auth_grpc_app",
		ExpiresAt: jwt.NewNumericDate(exp),
	}

	if jwtRefresh, err = JWT(u.cfg.JWTSecret(), claims); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if hash, err = u.hash(jwtRefresh); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return jwtRefresh, hash, nil
}

// hash generates a hash using sha256 and return base64 string
func (u *VerifyUseCase) hash(str string) (hash string, err error) {
	hasher := sha256.New()
	if _, err = hasher.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("hash: write failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}
