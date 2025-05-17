package authrepo

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/users/internal/core/repository/query"
)

var createSession = `
	insert into user_session (user_id, token) values($1, $2)
`

func (r *AuthRepository) CreateSession(ctx context.Context, userID int64, session string) error {
	const op = "AuthRepository.CreateSession"

	q := &query.Query{
		Name:      "Create user session",
		RawQuery:  createSession,
		Arguments: []any{userID, session},
	}

	if _, err := r.DB.Exec(ctx, q); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
