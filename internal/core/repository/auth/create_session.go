package authrepo

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/repository/query"
)

var createSession = `
	insert into refresh_tokens (jti, user_id, client_id, refresh_token_hash, expires_at) 
	values ($1, $2, $3, $4, $5) 
		on conflict (client_id) 
		    do update 
		    	set jti = $1, 
		    	    refresh_token_hash = $4, 
		    	    expires_at = $5;

`

func (r *AuthRepository) CreateSession(ctx context.Context, payload *entities.CreateSession) error {
	const op = "AuthRepository.CreateSession"

	q := &query.Query{
		Name:      "Create user session",
		RawQuery:  createSession,
		Arguments: []any{payload.RefreshTokenUUID, payload.UserID, payload.ClientUUID, payload.RefreshTokenHash, payload.ExpiresAt},
	}

	if _, err := r.DB.Exec(ctx, q); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
