package authrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

var refreshSession = `
with update_session as (
	update refresh_tokens
	set jti = $1
	where client_uuid = $2 and jti = $3 
	returning user_id
)
select 
	u.username,
	u.email,
	u.role_id 
from users u
where u.user_id = (select user_id from update_session);
`

func (r *AuthRepository) RefreshSession(ctx context.Context, payload *entities.RefreshSession) (*entities.TokenUserInfo, error) {
	const op = "AuthRepository.RefreshSession"

	q := &query.Query{
		Name:      "Refresh user session",
		RawQuery:  refreshSession,
		Arguments: []any{payload.NewRefreshTokenUUID, payload.ClientUUID, payload.OldRefreshTokenUUID},
	}

	userInfo := &entities.TokenUserInfo{}

	if err := r.DB.QueryRow(ctx, q).Scan(&userInfo.Username, &userInfo.Email, &userInfo.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrInvalidRefreshToken
		}
		return nil, apperrors.Wrap(op, err)
	}

	return userInfo, nil
}
