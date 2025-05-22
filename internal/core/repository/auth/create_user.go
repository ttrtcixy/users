package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

//var createUser = "insert into users (username, email, role_id) values ($1, $2, $3) returning user_id"
//var createUserPassword = "insert into user_password (user_id, password_hash, salt) "

var createUser = `
with new_user as (
	insert into users (username, email, role_id) 
	values ($1, $2, $3) 
	returning user_id
)
insert into user_password (user_id, hash, salt)
select 
	user_id, 
	$4, 
	$5 
from new_user;
`

// todo test

func (r *AuthRepository) CreateUser(ctx context.Context, payload *entities.CreateUserRequest) error {
	const op = "AuthRepository.CreateUser"

	q := &query.Query{
		Name:      "CreateUser",
		RawQuery:  createUser,
		Arguments: []any{payload.Username, payload.Email, payload.RoleID, payload.PasswordHash, payload.PasswordSalt},
	}

	_, err := r.DB.Exec(ctx, q)
	if err != nil {
		return apperrors.Wrap(op, err)
	}

	return nil
}
