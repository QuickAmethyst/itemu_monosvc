package errors

import "github.com/palantir/stacktrace"

const (
	ErrReadConfig = stacktrace.ErrorCode(iota)
	ErrUnmarshal
	ErrCodeUpgradeFailed
	ErrCodeReadyStateFiled
	ErrHttpServerListenerNil
)

var PropagateWithCode = stacktrace.PropagateWithCode
var Propagate = stacktrace.Propagate
