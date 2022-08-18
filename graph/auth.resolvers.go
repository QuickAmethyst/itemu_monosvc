package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/QuickAmethyst/monosvc/graph/model"
	accountUC "github.com/QuickAmethyst/monosvc/module/account/usecase"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
)

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.SignInInput) (*model.Credential, error) {
	accessTokenDetail, refreshTokenDetail, err := r.AccountUsecase.SignInWithEmail(ctx, accountUC.SignInWithEmailParams{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to sign in")
	}

	return &model.Credential{
		AccessToken:   accessTokenDetail.Raw,
		RefreshToken:  refreshTokenDetail.Raw,
		AccessExpire:  accessTokenDetail.Claims.ExpiresAt.Unix(),
		RefreshExpire: refreshTokenDetail.Claims.ExpiresAt.Unix(),
	}, nil
}

// RefreshCredential is the resolver for the refreshCredential field.
func (r *mutationResolver) RefreshCredential(ctx context.Context, input string) (*model.Credential, error) {
	accessTokenDetail, refreshTokenDetail, err := r.AccountUsecase.RefreshAccessToken(input)
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to refresh token")
	}

	return &model.Credential{
		AccessToken:   accessTokenDetail.Raw,
		RefreshToken:  refreshTokenDetail.Raw,
		AccessExpire:  accessTokenDetail.Claims.ExpiresAt.Unix(),
		RefreshExpire: refreshTokenDetail.Claims.ExpiresAt.Unix(),
	}, nil
}
