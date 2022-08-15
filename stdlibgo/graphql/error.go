package graphql

import "fmt"

type Error struct {
	error   error
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.error.Error())
}

func NewError(err error, message string) *Error {
	return &Error{err,message}
}
