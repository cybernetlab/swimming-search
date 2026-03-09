package domain

import (
	"errors"
	"fmt"
)

type ErrInvalidContext struct {
	Field string
}

type ErrAlreadyExists struct {
	Subject string
}

var (
	ErrNotFound       = errors.New("not found")
	ErrAPI            = errors.New("API error")
	ErrEmptyCentreIDs = errors.New("list of centre ID is empty")
)

func (e ErrInvalidContext) Error() string {
	return fmt.Sprintf("%s not found in context", e.Field)
}

func NewErrInvalidContext(field string) ErrInvalidContext {
	return ErrInvalidContext{Field: field}
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.Subject)
}

func NewErrAlreadyExists(subj string) ErrAlreadyExists {
	return ErrAlreadyExists{Subject: subj}
}
