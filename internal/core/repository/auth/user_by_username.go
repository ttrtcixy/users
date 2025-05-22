package authrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

var userByUsername = `
select 
    u.user_id, 
    u.email, 
    u.is_active, 
    u.role_id, 
    up.hash, 
    up.salt 
from users u
	join user_password up using(user_id) 
where username = $1;
`

func (r *AuthRepository) UserByUsername(ctx context.Context, username string) (*entities.User, error) {
	const op = "AuthRepository.UserByUsername"

	q := &query.Query{
		Name:      "Get user by username",
		RawQuery:  userByUsername,
		Arguments: []any{username},
	}
	var user = &entities.User{}

	err := r.DB.QueryRow(ctx, q).Scan(
		&user.ID,
		&user.Email,
		&user.IsActive,
		&user.RoleId,
		&user.PasswordHash,
		&user.PasswordSalt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotRegister
		}
		return nil, apperrors.Wrap(op, err)
	}

	user.Username = username

	return user, nil
}
