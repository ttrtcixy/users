package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

var checkLoginExists = `
select 
    EXISTS(SELECT 1 FROM users WHERE username = $1) AS username_exists, 
    EXISTS(SELECT 1 FROM users WHERE email = $2) AS email_exists;
`

// CheckLoginExist - check username and email are free
func (r *AuthRepository) CheckLoginExist(ctx context.Context, payload *entities.SignupRequest) (*entities.CheckLoginResponse, error) {
	const op = "AuthRepository.CheckLoginExist"

	q := &query.Query{
		Name:      "CheckLoginExists",
		RawQuery:  checkLoginExists,
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
		}, apperrors.Wrap(op, err)
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
