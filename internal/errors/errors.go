package apperrors

import (
	"fmt"
	"strings"
)

type ErrLoginExists struct {
	Username string
	Email    string
}

func (e *ErrLoginExists) Error() string {
	var str strings.Builder
	if e.Username != "" {
		str.WriteString(fmt.Sprintf("username: %s, exists; ", e.Username))
	}
	if e.Email != "" {
		str.WriteString(fmt.Sprintf("email: %s, exists; ", e.Email))
	}
	return strings.TrimSpace(str.String())
}
