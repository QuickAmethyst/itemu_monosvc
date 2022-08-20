package http

import (
	"github.com/QuickAmethyst/monosvc/stdlibgo/appcontext"
	"github.com/google/uuid"
	netHttp "net/http"
)

func appendHeaderToContext(handler netHttp.HandlerFunc) netHttp.HandlerFunc {
	return func(writer netHttp.ResponseWriter, request *netHttp.Request) {
		ctx := request.Context()

		requestID := request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = appcontext.SetRequestID(ctx, requestID)

		handler(writer, request.WithContext(ctx))
	}
}
