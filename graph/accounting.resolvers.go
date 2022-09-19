package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/QuickAmethyst/monosvc/graph/model"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
	libErr "github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
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
		TypeID:   uint(accountClass.TypeID),
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
		TypeID:   uint(accountClass.TypeID),
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

// AccountClasses is the resolver for the accountClasses field.
func (r *queryResolver) AccountClasses(ctx context.Context, input *model.AccountClassesInput) (*model.AccountClassesResult, error) {
	var (
		result model.AccountClassesResult
		p      qb.Paging
	)

	if input.Paging != nil {
		p.PageSize = input.Paging.PageSize
		p.CurrentPage = input.Paging.CurrentPage
	}

	accountClasses, paging, err := r.AccountingUsecase.GetAccountClassList(ctx, sql.AccountClassStatement{}, p)

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account classes", libErr.GetCode(err))
	}

	for _, accountClass := range accountClasses {
		result.Data = append(result.Data, &model.AccountClass{
			ID:       accountClass.ID,
			Name:     accountClass.Name,
			TypeID:   uint(accountClass.TypeID),
			Inactive: accountClass.Inactive,
		})
	}

	result.Paging = &model.Paging{
		CurrentPage: paging.CurrentPage,
		PageSize:    paging.PageSize,
		Total:       paging.Total,
	}

	return &result, nil
}

// AccountClass is the resolver for the accountClass field.
func (r *queryResolver) AccountClass(ctx context.Context, input model.AccountClassInput) (*model.AccountClass, error) {
	accountClass, err := r.AccountingUsecase.GetAccountClass(ctx, sql.AccountClassStatement{
		ID: int64(input.ID),
	})

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account class", libErr.GetCode(err))
	}

	return &model.AccountClass{
		ID:       accountClass.ID,
		Name:     accountClass.Name,
		TypeID:   uint(accountClass.TypeID),
		Inactive: accountClass.Inactive,
	}, nil
}
