package authrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/storage"
)

var req = "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) AS username_exists, EXISTS(SELECT 1 FROM users WHERE email = $2) AS email_exists;"

func (r *AuthRepository) CheckLoginExist(ctx context.Context, payload *entities.SignupRequest) (*entities.CheckLoginResponse, error) {
	query := storage.Query{
		QueryName: "CheckLoginExists",
		Query:     req,
		Args:      []any{payload.Username, payload.Email},
	}

	var (
		usernameExists bool
		emailExists    bool
	)

	err := r.DB.QueryRow(ctx, query).Scan(&usernameExists, &emailExists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entities.CheckLoginResponse{
				Status: false,
			}, nil
		}
		return &entities.CheckLoginResponse{
			Status: false,
		}, err
	}

	return &entities.CheckLoginResponse{
		Status:         true,
		UsernameExists: usernameExists,
		EmailExists:    emailExists,
	}, nil
}
