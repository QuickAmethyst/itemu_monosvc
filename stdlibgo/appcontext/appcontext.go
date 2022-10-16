package appcontext

import (
	"context"
	"github.com/google/uuid"
)

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

func GetUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(UserID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return v
}

func SetUserID(ctx context.Context, val uuid.UUID) context.Context {
	return context.WithValue(ctx, UserID, val)
}
