package apperrors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrServer            = errors.New("server error")
	ErrEmailTokenExpired = errors.New("email verification token is expired, request a new token")
	ErrEmailVerify       = errors.New("please verify email")
	ErrUserNotRegister   = errors.New("the user is not registered")
	ErrInvalidPassword   = errors.New("invalid password")
)

type UserError interface {
	UserError()
}

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
func (e *ErrLoginExists) UserError() {}

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}
func (e *ValidationError) UserError() {}

type ValidationErrors []ValidationError

func (e *ValidationErrors) UserError() {}
func (e *ValidationErrors) Add(field, err string) {
	*e = append(*e, ValidationError{
		Field: field,
		Err:   errors.New(err),
	})

}
func (e *ValidationErrors) Error() string {
	var err strings.Builder

	for _, v := range *e {
		err.WriteString(v.Field)
		err.WriteString(": ")
		err.WriteString(v.Err.Error() + "; ")
	}
	return strings.TrimSpace(err.String())
}
