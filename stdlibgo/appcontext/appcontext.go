package appcontext

import "context"

func GetRequestID(ctx context.Context) string {
	return ctx.Value(RequestIDKey).(string)
}

func SetRequestID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, RequestIDKey, val)
}

func GetBearerToken(ctx context.Context) string {
	v, ok := ctx.Value(BearerToken).(string)
	if !ok {
		return ""
	}

	return v
}

func SetBearerToken(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, BearerToken, val)
}

func GetUserID(ctx context.Context) string {
	v, ok := ctx.Value(UserID).(string)
	if !ok {
		return ""
	}

	return v
}

func SetUserID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, UserID, val)
}
