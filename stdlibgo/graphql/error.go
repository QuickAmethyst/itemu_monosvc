package graphql

import (
	"fmt"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
)

type Error struct {
	error   error
	Message string
	Code    errors.ErrorCode
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.error.Error())
}

func NewError(err error, message string, code errors.ErrorCode) *Error {
	return &Error{err, message, code}
}
