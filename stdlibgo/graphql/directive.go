package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/appcontext"
	"github.com/QuickAmethyst/monosvc/stdlibgo/auth"
)

type Directive func (ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error)

func AuthenticatedDirective(auth auth.Auth) Directive {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		bearer := appcontext.GetBearerToken(ctx)

		claim, err := auth.Authenticate(bearer)
		if err != nil {
			return nil, NewError(err, "authenticate failed")
		}

		ctx = appcontext.SetUserID(ctx, claim.Subject)

		return next(ctx)
	}
}
