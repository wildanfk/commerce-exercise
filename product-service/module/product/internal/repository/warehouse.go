package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"product-service/internal/util/liberr"
	"product-service/module/product/entity"
	"time"

	rh "github.com/hashicorp/go-retryablehttp"
)

type WarehouseConfiguration struct {
	ApiHost           string
	BasicAuthUsername string
	BasicAuthPassword string
}

type WarehouseRepository struct {
	Config     WarehouseConfiguration
	httpClient *rh.Client
}

func NewWarehouseRepository(config WarehouseConfiguration, httpClient *rh.Client) *WarehouseRepository {
	return &WarehouseRepository{
		Config:     config,
		httpClient: httpClient,
	}
}

type warehouseStock struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouse_id"`
	ProductID   string    `json:"product_id"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type warehouse struct {
	ID        string    `json:"id"`
	ShopID    string    `json:"shop_id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type listWarehouseStockResponse struct {
	Warehouses      []*warehouse      `json:"warehouses"`
	WarehouseStocks []*warehouseStock `json:"warehouse_stocks"`
	Meta            *entity.Meta      `json:"meta"`
}

func (w *WarehouseRepository) basicAuth() string {
	auth := w.Config.BasicAuthUsername + ":" + w.Config.BasicAuthPassword
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encodedAuth
}

func (w *WarehouseRepository) ActiveStock(ctx context.Context, productIDs []string) ([]*entity.WarehouseStock, error) {
	qparams := url.Values{}
	for _, pid := range productIDs {
		qparams.Add("product_ids", pid)
	}

	path := w.Config.ApiHost + "/active-stocks?" + qparams.Encode()

	req, err := rh.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", w.basicAuth())

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, liberr.NewTracer("Error when request on Warehouse.ActiveStock").Wrap(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, liberr.NewTracer("Error happened when read body response on Warehouse.ActiveStock").Wrap(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, liberr.NewTracer(fmt.Sprintf("Error with http status %d on Warehouse.ActiveStock", resp.StatusCode)).Wrap(err)
	}

	responseObj := listWarehouseStockResponse{}
	err = json.Unmarshal(responseBody, &responseObj)
	if err != nil {
		return nil, liberr.NewTracer("Error happened when parse json body response on Warehouse.ActiveStock").Wrap(err)
	}

	warehouseShopMap := map[string]*warehouse{}
	if responseObj.Warehouses != nil {
		for _, ws := range responseObj.Warehouses {
			warehouseShopMap[ws.ID] = ws
		}
	}

	warehouseStocks := []*entity.WarehouseStock{}
	if responseObj.WarehouseStocks != nil {
		for _, ws := range responseObj.WarehouseStocks {
			warehouseStocks = append(warehouseStocks, &entity.WarehouseStock{
				ID:            ws.ID,
				WarehouseID:   ws.WarehouseID,
				WarehouseName: warehouseShopMap[ws.WarehouseID].Name,
				ShopID:        warehouseShopMap[ws.WarehouseID].ShopID,
				ProductID:     ws.ProductID,
				Stock:         ws.Stock,
			})
		}
	}

	return warehouseStocks, nil
}
