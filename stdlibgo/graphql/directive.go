package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/appcontext"
	"github.com/QuickAmethyst/monosvc/stdlibgo/auth"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/google/uuid"
)

type Directive func (ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error)

func AuthenticatedDirective(auth auth.Auth) Directive {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		bearer := appcontext.GetBearerToken(ctx)

		claim, err := auth.Authenticate(bearer)
		if err != nil {
			return nil, NewError(err, "authenticate failed", errors.GetCode(err))
		}

		userID, err := uuid.Parse(claim.Subject)
		if err != nil {
			return nil, NewError(err, "authenticate failed", errors.GetCode(err))
		}

		ctx = appcontext.SetUserID(ctx, userID)

		return next(ctx)
	}
}
