package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/QuickAmethyst/monosvc/graph/model"
	inventorySql "github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
)

// StoreUom is the resolver for the storeUom field.
func (r *mutationResolver) StoreUom(ctx context.Context, input model.WriteUomInput) (*model.Uom, error) {
	uom, err := input.Domain()
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to read input")
	}

	if err := r.Resolver.InventoryUsecase.StoreUom(ctx, &uom); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on create uom")
	}

	return &model.Uom{
		ID:          uom.ID,
		Name:        uom.Name,
		Description: uom.Description.String,
		Decimal:     uom.Decimal.Int32,
	}, nil
}

// UpdateUom is the resolver for the updateUom field.
func (r *mutationResolver) UpdateUom(ctx context.Context, id int, input model.WriteUomInput) (*model.Uom, error) {
	uom, err := input.Domain()
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to read input")
	}

	if err = r.Resolver.InventoryUsecase.UpdateUomByID(ctx, int64(id), &uom); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on update uom")
	}

	return &model.Uom{
		ID:          int64(id),
		Name:        uom.Name,
		Description: uom.Description.String,
		Decimal:     uom.Decimal.Int32,
	}, nil
}

// Uoms is the resolver for the uoms field.
func (r *queryResolver) Uoms(ctx context.Context, input *model.UomsInput) (*model.UomsResult, error) {
	var (
		result model.UomsResult
		p      qb.Paging
	)

	p.Normalize()
	if input != nil {
		if input.Paging != nil {
			p.CurrentPage = input.Paging.CurrentPage
			p.PageSize = input.Paging.PageSize
		}
	}

	uoms, paging, err := r.InventoryUsecase.GetUomList(ctx, inventorySql.UomStatement{}, p)

	if err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed on get list of uom")
	}

	for _, uom := range uoms {
		result.Data = append(result.Data, &model.Uom{
			ID:          uom.ID,
			Name:        uom.Name,
			Description: uom.Description.String,
			Decimal:     uom.Decimal.Int32,
		})
	}

	result.Paging = &model.Paging{
		CurrentPage: paging.CurrentPage,
		PageSize:    paging.PageSize,
		Total:       paging.Total,
	}

	return &result, nil
}
