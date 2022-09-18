package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/QuickAmethyst/monosvc/graph/model"
	libErr "github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
)

// StoreAccountClass is the resolver for the storeAccountClass field.
func (r *mutationResolver) StoreAccountClass(ctx context.Context, input model.WriteAccountClassesInput) (*model.AccountClass, error) {
	accountClass, err := input.Domain()
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to read input", libErr.GetCode(err))
	}

	if err = r.Resolver.AccountingUsecase.StoreAccountClass(ctx, &accountClass); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on create account class", libErr.GetCode(err))
	}

	return &model.AccountClass{
		ID:       accountClass.ID,
		Name:     accountClass.Name,
		Type:     uint(accountClass.Type),
		Inactive: accountClass.Inactive,
	}, nil
}

// UpdateAccountClassByID is the resolver for the updateAccountClassByID field.
func (r *mutationResolver) UpdateAccountClassByID(ctx context.Context, id int, input model.WriteAccountClassesInput) (*model.AccountClass, error) {
	accountClass, err := input.Domain()
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to read input", libErr.GetCode(err))
	}

	if err = r.Resolver.AccountingUsecase.UpdateAccountClassByID(ctx, int64(id), &accountClass); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on update account class", libErr.GetCode(err))
	}

	return &model.AccountClass{
		ID:       int64(id),
		Name:     accountClass.Name,
		Type:     uint(accountClass.Type),
		Inactive: accountClass.Inactive,
	}, nil
}

// DeleteAccountClassByID is the resolver for the deleteAccountClassByID field.
func (r *mutationResolver) DeleteAccountClassByID(ctx context.Context, id int) (int, error) {
	if err := r.AccountingUsecase.DeleteAccountClassByID(ctx, int64(id)); err != nil {
		r.Logger.Error(err.Error())
		return id, sdkGraphql.NewError(err, "Failed to delete account class", libErr.GetCode(err))
	}

	return id, nil
}
