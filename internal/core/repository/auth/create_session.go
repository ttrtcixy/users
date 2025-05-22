package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

//WITH inserted AS (
//INSERT INTO refresh_tokens (jti, user_id, client_uuid, refresh_token_hash, expires_at)
//VALUES ($1, $2, $3, $4, $5)
//ON CONFLICT (client_uuid) DO UPDATE
//SET jti = EXCLUDED.jti,
//refresh_token_hash = EXCLUDED.refresh_token_hash,
//expires_at = EXCLUDED.expires_at
//RETURNING user_id
//)
//DELETE FROM refresh_tokens
//WHERE id IN (
//SELECT id FROM refresh_tokens
//WHERE user_id = (SELECT user_id FROM inserted)
//ORDER BY created_at ASC
//OFFSET 5
//);

// todo add max session count.
var createSession = `
insert into refresh_tokens (jti, user_id, client_uuid, expires_at) 
values ($1, $2, $3, $4) 
	on conflict (client_uuid) 
	do update 
		set jti = $1,  
		    expires_at = $4;
`

// CreateSession - create new session to user with clientID and token jti
func (r *AuthRepository) CreateSession(ctx context.Context, payload *entities.CreateSession) error {
	const op = "AuthRepository.CreateSession"

	q := &query.Query{
		Name:      "Create user session",
		RawQuery:  createSession,
		Arguments: []any{payload.RefreshTokenUUID, payload.UserID, payload.ClientUUID, payload.ExpiresAt},
	}

	if _, err := r.DB.Exec(ctx, q); err != nil {
		return apperrors.Wrap(op, err)
	}

	return nil
}
