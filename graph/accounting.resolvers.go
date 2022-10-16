package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/QuickAmethyst/monosvc/graph/generated"
	"github.com/QuickAmethyst/monosvc/graph/model"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/appcontext"
	libErr "github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
)

// Group is the resolver for the group field.
func (r *accountResolver) Group(ctx context.Context, obj *model.Account) (*model.AccountGroup, error) {
	if obj == nil || obj.GroupID == 0 {
		return nil, nil
	}

	accountGroup, err := r.AccountingUsecase.GetAccountGroup(ctx, sql.AccountGroupStatement{ID: obj.GroupID})
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account group", libErr.GetCode(err))
	}

	return &model.AccountGroup{
		ID:       accountGroup.ID,
		Name:     accountGroup.Name,
		ClassID:  accountGroup.ClassID,
		ParentID: accountGroup.ParentID.Int64,
		Inactive: accountGroup.Inactive,
	}, nil
}

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
		ID:       int64(id),
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

// StoreAccount is the resolver for the storeAccount field.
func (r *mutationResolver) StoreAccount(ctx context.Context, input model.WriteAccountInput) (*model.Account, error) {
	account := input.Domain()

	if err := r.AccountingUsecase.StoreAccount(ctx, &account); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on store account", libErr.GetCode(err))
	}

	return &model.Account{
		ID:       account.ID,
		Name:     account.Name,
		GroupID:  account.GroupID,
		Inactive: account.Inactive,
	}, nil
}

// UpdateAccountByID is the resolver for the updateAccountByID field.
func (r *mutationResolver) UpdateAccountByID(ctx context.Context, id int, input model.WriteAccountInput) (*model.Account, error) {
	account := input.Domain()

	if err := r.AccountingUsecase.UpdateAccountByID(ctx, int64(id), &account); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on update account by id", libErr.GetCode(err))
	}

	return &model.Account{
		ID:       int64(id),
		Name:     account.Name,
		GroupID:  account.GroupID,
		Inactive: account.Inactive,
	}, nil
}

// DeleteAccountByID is the resolver for the deleteAccountByID field.
func (r *mutationResolver) DeleteAccountByID(ctx context.Context, id int) (int, error) {
	if err := r.AccountingUsecase.DeleteAccountByID(ctx, int64(id)); err != nil {
		r.Logger.Error(err.Error())
		return id, sdkGraphql.NewError(err, "Failed on delete account by id", libErr.GetCode(err))
	}

	return id, nil
}

// StoreTransactions is the resolver for the storeTransactions field.
func (r *mutationResolver) StoreTransactions(ctx context.Context, input []*model.WriteTransactionsInput) (*model.Journal, error) {
	userID := appcontext.GetUserID(ctx)

	transactions := make([]sql.Transaction, len(input))
	for _, item := range input {
		transactions = append(transactions, sql.Transaction{
			AccountID: item.AccountID,
			Amount:    item.Amount,
		})
	}

	journal, err := r.AccountingUsecase.StoreTransactions(ctx, userID, transactions)
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on create transactions", libErr.GetCode(err))
	}

	if journal == nil {
		return nil, nil
	}

	return &model.Journal{
		ID:        journal.ID.String(),
		Amount:    journal.Amount,
		CreatedAt: journal.CreatedAt,
	}, nil
}

// AccountClasses is the resolver for the accountClasses field.
func (r *queryResolver) AccountClasses(ctx context.Context) ([]*model.AccountClass, error) {
	accountClasses, err := r.AccountingUsecase.GetAllAccountClasses(ctx, sql.AccountClassStatement{})
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
	classTypes := r.AccountingUsecase.GetAllAccountTypes(ctx)

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

	accountGroups, err := r.AccountingUsecase.GetAllAccountGroups(ctx, statement)

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

// Accounts is the resolver for the accounts field.
func (r *queryResolver) Accounts(ctx context.Context, input *model.AccountInput) ([]*model.Account, error) {
	var stmt sql.AccountStatement
	if input != nil {
		stmt.ID = input.ID
	}

	accounts, err := r.AccountingUsecase.GetAllAccounts(ctx, stmt)
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get accounts", libErr.GetCode(err))
	}

	var result = make([]*model.Account, len(accounts))
	for i, account := range accounts {
		result[i] = &model.Account{
			ID:       account.ID,
			Name:     account.Name,
			GroupID:  account.GroupID,
			Inactive: account.Inactive,
		}
	}

	return result, nil
}

// Account is the resolver for the account field.
func (r *queryResolver) Account(ctx context.Context, input model.AccountInput) (*model.Account, error) {
	account, err := r.AccountingUsecase.GetAccount(ctx, sql.AccountStatement{
		ID: input.ID,
	})

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get account group", libErr.GetCode(err))
	}

	return &model.Account{
		ID:       account.ID,
		Name:     account.Name,
		GroupID:  account.GroupID,
		Inactive: account.Inactive,
	}, nil
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

// AccountClass returns generated.AccountClassResolver implementation.
func (r *Resolver) AccountClass() generated.AccountClassResolver { return &accountClassResolver{r} }

// AccountGroup returns generated.AccountGroupResolver implementation.
func (r *Resolver) AccountGroup() generated.AccountGroupResolver { return &accountGroupResolver{r} }

type accountResolver struct{ *Resolver }
type accountClassResolver struct{ *Resolver }
type accountGroupResolver struct{ *Resolver }
