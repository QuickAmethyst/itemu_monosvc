package errors

import "github.com/palantir/stacktrace"

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
