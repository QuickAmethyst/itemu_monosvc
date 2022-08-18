package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/QuickAmethyst/monosvc/graph/model"
	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
	inventorySql "github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"
	sdkGraphql "github.com/QuickAmethyst/monosvc/stdlibgo/graphql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
)

// StoreUom is the resolver for the storeUom field.
func (r *mutationResolver) StoreUom(ctx context.Context, input model.WriteUomInput) (*model.Uom, error) {
	var uom domain.Uom
	uom.Name = input.Name

	if err := uom.Description.Scan(input.Description); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to assign description")
	}

	if err := uom.Decimal.Scan(input.Decimal); err != nil {
		r.Logger.Error(err.Error())
		return nil, sdkGraphql.NewError(err, "Failed to assign decimal")
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
