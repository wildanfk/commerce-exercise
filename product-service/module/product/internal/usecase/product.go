package usecase

import (
	"context"
	"product-service/internal/util/liberr"
	"product-service/internal/util/libpagination"
	"product-service/module/product/entity"
)

type ProductUsecaseRepos struct {
	ProductRepo   ProductRepository
	WarehouseRepo WarehouseRepository
	ShopRepo      ShopRepository
}

type ProductUsecase struct {
	repos *ProductUsecaseRepos
}

func NewProductUsecase(repos *ProductUsecaseRepos) *ProductUsecase {
	return &ProductUsecase{
		repos: repos,
	}
}

func (p *ProductUsecase) CheckProduct(ctx context.Context, params *entity.ListProductByParams) ([]*entity.Product, *libpagination.OffsetPagination, error) {
	params.Offset = libpagination.Offset(params.Page, params.Limit)

	return p.repos.ProductRepo.ListByParams(ctx, params)
}

func (p *ProductUsecase) ListProduct(ctx context.Context, params *entity.ListProductByParams) ([]*entity.ProductDetail, *libpagination.OffsetPagination, error) {
	params.Offset = libpagination.Offset(params.Page, params.Limit)

	// Retrieve products
	products, pagination, err := p.repos.ProductRepo.ListByParams(ctx, params)
	if err != nil {
		return nil, nil, liberr.ResolveError(err)
	}

	productDetails := []*entity.ProductDetail{}
	if len(products) == 0 {
		return productDetails, pagination, nil
	}

	productIDs := []string{}
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}

	// Retrieve warehouse stocks
	warehouseStocks, err := p.repos.WarehouseRepo.ActiveStock(ctx, productIDs)
	if err != nil {
		return nil, nil, liberr.ResolveError(err)
	}

	// Get Unique shop id
	shopIDsMap := make(map[string]struct{})
	shopIDs := []string{}
	for _, sa := range warehouseStocks {
		if _, exists := shopIDsMap[sa.ShopID]; !exists {
			shopIDsMap[sa.ShopID] = struct{}{}
			shopIDs = append(shopIDs, sa.ShopID)
		}
	}

	// Retrieve shops
	shops := []*entity.Shop{}
	if len(shopIDs) > 0 {
		shops, err = p.repos.ShopRepo.ListByShopIDs(ctx, shopIDs)
		if err != nil {
			return nil, nil, liberr.ResolveError(err)
		}
	}

	// Build shop name map
	shopNameMap := map[string]string{}
	for _, s := range shops {
		shopNameMap[s.ID] = s.Name
	}

	// map[product_id][shop_id]detailWarehouse
	productShopWarehouseMap := map[string]map[string][]*entity.ProductDetailWarehouse{}

	// map[product_id][shop_id]total_stock
	productShopTotalStockMap := map[string]map[string]int{}

	// map[product_id]total_stock
	productTotalStockMap := map[string]int{}

	// Build product shop warehouse
	for _, ws := range warehouseStocks {
		if _, ok := productShopTotalStockMap[ws.ProductID]; !ok {
			productShopWarehouseMap[ws.ProductID] = map[string][]*entity.ProductDetailWarehouse{}
			productShopTotalStockMap[ws.ProductID] = map[string]int{}
			productTotalStockMap[ws.ProductID] = 0
		}

		if _, ok := productShopTotalStockMap[ws.ProductID][ws.ShopID]; !ok {
			productShopWarehouseMap[ws.ProductID][ws.ShopID] = []*entity.ProductDetailWarehouse{}
			productShopTotalStockMap[ws.ProductID][ws.ShopID] = 0
		}

		productShopWarehouseMap[ws.ProductID][ws.ShopID] = append(productShopWarehouseMap[ws.ProductID][ws.ShopID], &entity.ProductDetailWarehouse{
			WarehouseID:      ws.WarehouseID,
			WarehouseName:    ws.WarehouseName,
			WarehouseStockID: ws.ID,
			WarehouseStock:   ws.Stock,
		})
		productShopTotalStockMap[ws.ProductID][ws.ShopID] += ws.Stock
		productTotalStockMap[ws.ProductID] += ws.Stock
	}

	// Build Product Detail
	for _, p := range products {
		totalStock, _ := productTotalStockMap[p.ID]

		productDetailShops := []*entity.ProductDetailShop{}
		for _, sid := range shopIDs {
			if ws, ok := productShopWarehouseMap[p.ID][sid]; ok {
				productDetailShops = append(productDetailShops, &entity.ProductDetailShop{
					ID:         sid,
					Name:       shopNameMap[sid],
					TotalStock: productShopTotalStockMap[p.ID][sid],
					Warehouses: ws,
				})
			}
		}

		productDetails = append(productDetails, &entity.ProductDetail{
			ID:         p.ID,
			Name:       p.Name,
			Price:      p.Price,
			TotalStock: totalStock,
			Shops:      productDetailShops,
		})
	}

	return productDetails, pagination, nil
}
