package authrepo

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
)

var req = "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) AS username_exists, EXISTS(SELECT 1 FROM users WHERE email = $2) AS email_exists;"

func (r *AuthRepository) CheckLoginExist(ctx context.Context, payload *entities.SignupRequest) (*entities.CheckLoginResponse, error) {
	const op = "AuthRepository.CheckLoginExist"

	q := &query.Query{
		Name:      "CheckLoginExists",
		RawQuery:  req,
		Arguments: []any{payload.Username, payload.Email},
	}

	var (
		usernameExists bool
		emailExists    bool
	)

	err := r.DB.QueryRow(ctx, q).Scan(&usernameExists, &emailExists)
	if err != nil {
		return &entities.CheckLoginResponse{
			Status: false,
		}, fmt.Errorf("%s: %w", op, err)
	}

	if usernameExists == false && emailExists == false {
		return &entities.CheckLoginResponse{
			Status: false,
		}, nil
	}

	return &entities.CheckLoginResponse{
		Status:         true,
		UsernameExists: usernameExists,
		EmailExists:    emailExists,
	}, nil
}
