package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/QuickAmethyst/monosvc/graph/generated"
	"github.com/QuickAmethyst/monosvc/graph/model"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
	libErr "github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
)

// Type is the resolver for the type field.
func (r *accountClassResolver) Type(ctx context.Context, obj *model.AccountClass) (*model.AccountClassType, error) {
	accountClassType := r.AccountingUsecase.GetAccountClassTypeByID(ctx, obj.TypeID)
	return &model.AccountClassType{
		ID:   accountClassType.ID,
		Name: accountClassType.Name,
	}, nil
}

// Parent is the resolver for the parent field.
func (r *accountGroupResolver) Parent(ctx context.Context, obj *model.AccountGroup) (*model.AccountGroup, error) {
	if obj == nil || obj.ParentID == 0 {
		return nil, nil
	}

	accountGroup, err := r.AccountingUsecase.GetAccountGroupByID(ctx, obj.ParentID)
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account group parent", libErr.GetCode(err))
	}

	return &model.AccountGroup{
		ID:       accountGroup.ID,
		Name:     accountGroup.Name,
		ClassID:  accountGroup.ClassID,
		ParentID: accountGroup.ParentID.Int64,
		Inactive: accountGroup.Inactive,
	}, nil
}

// Class is the resolver for the class field.
func (r *accountGroupResolver) Class(ctx context.Context, obj *model.AccountGroup) (*model.AccountClass, error) {
	if obj == nil || obj.ClassID == 0 {
		return nil, nil
	}

	accountClass, err := r.AccountingUsecase.GetAccountClassByID(ctx, obj.ClassID)
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account group class", libErr.GetCode(err))
	}

	return &model.AccountClass{
		ID:       accountClass.ID,
		Name:     accountClass.Name,
		TypeID:   accountClass.TypeID,
		Inactive: accountClass.Inactive,
	}, nil
}

// Child is the resolver for the child field.
func (r *accountGroupResolver) Child(ctx context.Context, obj *model.AccountGroup) ([]*model.AccountGroup, error) {
	accountGroups, err := r.AccountingUsecase.GetAllTopLevelAccountGroup(ctx, sql.AccountGroupStatement{
		ParentID: obj.ID,
	})

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account group childs", libErr.GetCode(err))
	}

	result := make([]*model.AccountGroup, len(accountGroups))
	for i, accountGroup := range accountGroups {
		result[i] = &model.AccountGroup{
			ID:       accountGroup.ID,
			Name:     accountGroup.Name,
			ClassID:  accountGroup.ClassID,
			ParentID: accountGroup.ParentID.Int64,
			Inactive: false,
		}
	}

	return result, nil
}

// StoreAccountClass is the resolver for the storeAccountClass field.
func (r *mutationResolver) StoreAccountClass(ctx context.Context, input model.WriteAccountClassInput) (*model.AccountClass, error) {
	accountClass := input.Domain()

	if err := r.Resolver.AccountingUsecase.StoreAccountClass(ctx, &accountClass); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on create account class", libErr.GetCode(err))
	}

	return &model.AccountClass{
		ID:       accountClass.ID,
		Name:     accountClass.Name,
		TypeID:   accountClass.TypeID,
		Inactive: accountClass.Inactive,
	}, nil
}

// UpdateAccountClassByID is the resolver for the updateAccountClassByID field.
func (r *mutationResolver) UpdateAccountClassByID(ctx context.Context, id int, input model.WriteAccountClassInput) (*model.AccountClass, error) {
	accountClass := input.Domain()

	if err := r.Resolver.AccountingUsecase.UpdateAccountClassByID(ctx, int64(id), &accountClass); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on update account class", libErr.GetCode(err))
	}

	return &model.AccountClass{
		ID:       int64(id),
		Name:     accountClass.Name,
		TypeID:   accountClass.TypeID,
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

// StoreAccountGroup is the resolver for the storeAccountGroup field.
func (r *mutationResolver) StoreAccountGroup(ctx context.Context, input model.WriteAccountGroupInput) (*model.AccountGroup, error) {
	accountGroup, err := input.Domain()
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to create account group", libErr.GetCode(err))
	}

	if err = r.Resolver.AccountingUsecase.StoreAccountGroup(ctx, &accountGroup); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to create account group", libErr.GetCode(err))
	}

	return &model.AccountGroup{
		ID:       accountGroup.ID,
		Name:     accountGroup.Name,
		ClassID:  accountGroup.ClassID,
		ParentID: accountGroup.ParentID.Int64,
		Inactive: accountGroup.Inactive,
	}, nil
}

// UpdateAccountGroupByID is the resolver for the updateAccountGroupByID field.
func (r *mutationResolver) UpdateAccountGroupByID(ctx context.Context, id int, input model.WriteAccountGroupInput) (*model.AccountGroup, error) {
	accountGroup, err := input.Domain()
	if err != nil {
		return nil, sdkGraphql.NewError(err, "Failed to update account group", libErr.GetCode(err))
	}

	if err := r.Resolver.AccountingUsecase.UpdateAccountGroupByID(ctx, int64(id), &accountGroup); err != nil {
		return nil, sdkGraphql.NewError(err, "Failed to update account group", libErr.GetCode(err))
	}

	return &model.AccountGroup{
		ID:       accountGroup.ID,
		Name:     accountGroup.Name,
		ClassID:  accountGroup.ClassID,
		ParentID: accountGroup.ParentID.Int64,
		Inactive: accountGroup.Inactive,
	}, nil
}

// DeleteAccountGroupByID is the resolver for the deleteAccountGroupByID field.
func (r *mutationResolver) DeleteAccountGroupByID(ctx context.Context, id int) (int, error) {
	if err := r.AccountingUsecase.DeleteAccountGroupByID(ctx, int64(id)); err != nil {
		r.Logger.Error(err.Error())
		return id, sdkGraphql.NewError(err, "Failed to delete account group", libErr.GetCode(err))
	}

	return id, nil
}

// AccountClasses is the resolver for the accountClasses field.
func (r *queryResolver) AccountClasses(ctx context.Context) ([]*model.AccountClass, error) {
	accountClasses, err := r.AccountingUsecase.GetAllAccountClass(ctx, sql.AccountClassStatement{})
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account classes", libErr.GetCode(err))
	}

	result := make([]*model.AccountClass, len(accountClasses))
	for i, accountClass := range accountClasses {
		result[i] = &model.AccountClass{
			ID:       accountClass.ID,
			Name:     accountClass.Name,
			TypeID:   accountClass.TypeID,
			Inactive: accountClass.Inactive,
		}
	}

	return result, nil
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
		TypeID:   accountClass.TypeID,
		Inactive: accountClass.Inactive,
	}, nil
}

// AccountClassTypes is the resolver for the accountClassTypes field.
func (r *queryResolver) AccountClassTypes(ctx context.Context) (*model.AccountClassTypesResult, error) {
	result := make([]*model.AccountClassType, 0)
	classTypes := r.AccountingUsecase.GetAccountClassTypeList(ctx)

	for _, classType := range classTypes {
		result = append(result, &model.AccountClassType{
			ID:   classType.ID,
			Name: classType.Name,
		})
	}

	return &model.AccountClassTypesResult{Data: result}, nil
}

// AccountClassType is the resolver for the accountClassType field.
func (r *queryResolver) AccountClassType(ctx context.Context, input model.AccountClassTypeInput) (*model.AccountClassType, error) {
	classType := r.AccountingUsecase.GetAccountClassTypeByID(ctx, input.ID)
	return &model.AccountClassType{
		ID:   classType.ID,
		Name: classType.Name,
	}, nil
}

// AccountGroups is the resolver for the accountGroups field.
func (r *queryResolver) AccountGroups(ctx context.Context, input *model.AccountGroupInput) ([]*model.AccountGroup, error) {
	var statement sql.AccountGroupStatement
	if input != nil {
		statement = sql.AccountGroupStatement{ParentIDIsNULL: input.ParentIDIsNULL}
	}

	accountGroups, err := r.AccountingUsecase.GetAllAccountGroup(ctx, statement)

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account groups", libErr.GetCode(err))
	}

	result := make([]*model.AccountGroup, len(accountGroups))
	for i, accountGroup := range accountGroups {
		result[i] = &model.AccountGroup{
			ID:       accountGroup.ID,
			Name:     accountGroup.Name,
			ClassID:  accountGroup.ClassID,
			ParentID: accountGroup.ParentID.Int64,
			Inactive: accountGroup.Inactive,
		}
	}

	return result, nil
}

// AccountGroup is the resolver for the accountGroup field.
func (r *queryResolver) AccountGroup(ctx context.Context, input model.AccountGroupInput) (*model.AccountGroup, error) {
	accountGroup, err := r.AccountingUsecase.GetAccountGroup(ctx, sql.AccountGroupStatement{
		ID: int64(input.ID),
	})

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account group", libErr.GetCode(err))
	}

	return &model.AccountGroup{
		ID:       accountGroup.ID,
		Name:     accountGroup.Name,
		ParentID: accountGroup.ParentID.Int64,
		ClassID:  accountGroup.ClassID,
		Inactive: accountGroup.Inactive,
	}, nil
}

// AccountClass returns generated.AccountClassResolver implementation.
func (r *Resolver) AccountClass() generated.AccountClassResolver { return &accountClassResolver{r} }

// AccountGroup returns generated.AccountGroupResolver implementation.
func (r *Resolver) AccountGroup() generated.AccountGroupResolver { return &accountGroupResolver{r} }

type accountClassResolver struct{ *Resolver }
type accountGroupResolver struct{ *Resolver }
