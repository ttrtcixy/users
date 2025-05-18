package authrepo

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
)

//var createUser = "insert into users (username, email, role_id) values ($1, $2, $3) returning user_id"
//var createUserPassword = "insert into user_password (user_id, password_hash, salt) "

var createUser = `
with addUser as (
	insert into users (username, email, role_id) values ($1, $2, $3) returning user_id
)
insert into user_password (user_id, hash, salt)
	select user_id, $4, $5 from addUser;
`

func (r *AuthRepository) CreateUser(ctx context.Context, payload *entities.CreateUserRequest) error {
	const op = "AuthRepository.CreateUser"

	q := &query.Query{
		Name:      "CreateUser",
		RawQuery:  createUser,
		Arguments: []any{payload.Username, payload.Email, 2, payload.PasswordHash, payload.PasswordSalt},
	}

	_, err := r.DB.Exec(ctx, q)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
