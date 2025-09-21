package handler

import (
	"encoding/json"
	"net/http"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/internal/util/librest"
	"warehouse-service/module/warehouse/entity"
)

type WarehouseStockHandler struct {
	warehouseStockUsecase WarehouseStockUsecase
}

func NewWarehouseStockHandler(warehouseStockUsecase WarehouseStockUsecase) *WarehouseStockHandler {
	return &WarehouseStockHandler{
		warehouseStockUsecase: warehouseStockUsecase,
	}
}

func (ws *WarehouseStockHandler) ActiveStock(w http.ResponseWriter, r *http.Request) error {
	// Query parameters
	qparams := r.URL.Query()

	params := &entity.ListWarehouseStockByParams{
		ProductIDs: qparams["product_ids"],
	}

	warehouses, warehouseStocks, err := ws.warehouseStockUsecase.ActiveStock(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.ListWarehouseStockResponse{
		Warehouse:       warehouses,
		WarehouseStocks: warehouseStocks,
		Meta: &entity.Meta{
			HttpStatusCode: code,
		},
	}, code)
	return nil
}

func (ws *WarehouseStockHandler) AdjustmentStock(w http.ResponseWriter, r *http.Request) error {
	params := new(entity.WarehouseStockAdjustmentRequest)
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		return liberr.NewBaseError(entity.ErrorInvalidBodyJSON)
	}

	err := ws.warehouseStockUsecase.AdjustmentStock(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.GetMessageResponse{
		Message: "Success adjustment stock",
		Meta: &entity.Meta{
			HttpStatusCode: code,
		},
	}, code)
	return nil
}

func (ws *WarehouseStockHandler) TransferStock(w http.ResponseWriter, r *http.Request) error {
	params := new(entity.WarehouseStockTransferRequest)
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		return liberr.NewBaseError(entity.ErrorInvalidBodyJSON)
	}

	err := ws.warehouseStockUsecase.TransferStock(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.GetMessageResponse{
		Message: "Success transfer stock",
		Meta: &entity.Meta{
			HttpStatusCode: code,
		},
	}, code)
	return nil
}
