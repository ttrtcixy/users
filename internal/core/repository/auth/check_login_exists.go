package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
)

var req = "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) AS username_exists, EXISTS(SELECT 1 FROM users WHERE email = $2) AS email_exists;"

func (r *AuthRepository) CheckLoginExist(ctx context.Context, payload *entities.SignupRequest) (*entities.CheckLoginResponse, error) {
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
		}, err
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
