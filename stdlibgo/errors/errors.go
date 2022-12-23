package errors

import (
	"bytes"
	"fmt"
	"github.com/palantir/stacktrace"
	"strings"
)

type ErrorCode = stacktrace.ErrorCode

const (
	ErrReadConfig = stacktrace.ErrorCode(iota)
	ErrUnmarshal
	ErrCodeUpgradeFailed
	ErrCodeReadyStateFiled
	ErrHttpServerListenerNil
)

var PropagateWithCode = stacktrace.PropagateWithCode
var Propagate = stacktrace.Propagate
var GetCode = stacktrace.GetCode
var RootCause = stacktrace.RootCause

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (f *FieldError) Error() string {
	return fmt.Sprintf("Field %s. Error: %s", f.Field, f.Message)
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(v); i++ {
		fe := v[i]
		buff.WriteString(fe.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
