package graphql

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"
	netHttp "net/http"
)

type Options struct {
	Development bool
}

func New(es graphql.ExecutableSchema, opt Options) (h netHttp.HandlerFunc, pH netHttp.HandlerFunc) {
	srv := handler.NewDefaultServer(es)
	playgroundH := playground.Handler("GraphQL playground", "/graphql/query")

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		if !opt.Development {
			graphql.GetOperationContext(ctx).DisableIntrospection = true
		}

		return next(ctx)
	})

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		var graphErr *Error
		if errors.As(e, &graphErr) {
			err.Message = graphErr.Message
			err.Extensions = make(map[string]interface{})
			err.Extensions["code"] = graphErr.Code
			if opt.Development {
				err.Extensions["error"] = graphErr.Error()
			}
		}

		return err
	})

	return func(w netHttp.ResponseWriter, r *netHttp.Request) {
			srv.ServeHTTP(w, r)
		}, func(writer netHttp.ResponseWriter, request *netHttp.Request) {
			playgroundH.ServeHTTP(writer, request)
		}
}
