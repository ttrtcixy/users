package authrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

var activateUser = `
update users
	set is_active = true
where email = $1 and is_active = false
returning user_id, username, role_id;
`

// ActivateUser - if the user is not activated, activate them
func (r *AuthRepository) ActivateUser(ctx context.Context, email string) (*entities.TokenUserInfo, error) {
	const op = "AuthRepository.ActivateUser"
	q := &query.Query{
		Name:      "Activate User",
		RawQuery:  activateUser,
		Arguments: []any{email},
	}

	user := &entities.TokenUserInfo{}
	user.Email = email

	if err := r.DB.QueryRow(ctx, q).Scan(&user.ID, &user.Username, &user.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserAlreadyActivated
		}
		return nil, apperrors.Wrap(op, err)
	}

	return user, nil
}
