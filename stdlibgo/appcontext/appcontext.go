package appcontext

import "context"

func GetRequestID(ctx context.Context) string {
	return ctx.Value(RequestIDKey).(string)
}

func SetRequestID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, RequestIDKey, val)
}

func GetBearerToken(ctx context.Context) string {
	return ctx.Value(BearerToken).(string)
}

func SetBearerToken(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, BearerToken, val)
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(UserID).(string)
}

func SetUserID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, UserID, val)
}
