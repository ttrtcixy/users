package authrepo

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
)

// todo добавить проверку на то что пользователь уже активирован?
var activateUser = `
	update users
	set is_active = true
	where email = $1 returning user_id, username, role_id;
`

func (r *AuthRepository) ActivateUser(ctx context.Context, email string) (*entities.User, error) {
	const op = "AuthRepository.ActivateUser"
	q := &query.Query{
		Name:      "Activate User",
		RawQuery:  activateUser,
		Arguments: []any{email},
	}

	user := &entities.User{}
	user.Email = email

	if err := r.DB.QueryRow(ctx, q).Scan(&user.ID, &user.Username, &user.RoleId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
