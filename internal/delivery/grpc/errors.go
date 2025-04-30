package grpc

import (
	"errors"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type ValidationErrors []ValidationError

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
		err.WriteString(v.Err.Error())
		err.WriteString("\n")
	}
	return err.String()
}
