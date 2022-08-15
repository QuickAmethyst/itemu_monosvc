package graph

import (
	inventoryUC "github.com/QuickAmethyst/monosvc/module/inventory/usecase"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger           logger.Logger
	InventoryUsecase inventoryUC.Usecase
}
