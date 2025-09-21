package handler

import (
	"encoding/json"
	"net/http"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/internal/util/librest"
	"warehouse-service/module/warehouse/entity"
)

type WarehouseHandler struct {
	warehouseUsecase WarehouseUsecase
}

func NewWarehouseHandler(warehouseUsecase WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseUsecase: warehouseUsecase,
	}
}

func (ws *WarehouseHandler) WarehouseActivation(w http.ResponseWriter, r *http.Request) error {
	params := new(entity.WarehouseActivationRequest)
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		return liberr.NewBaseError(entity.ErrorInvalidBodyJSON)
	}

	err := ws.warehouseUsecase.WarehouseActivation(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.GetMessageResponse{
		Message: "Success inactive/active warehouse",
		Meta: &entity.Meta{
			HttpStatusCode: code,
		},
	}, code)
	return nil
}
