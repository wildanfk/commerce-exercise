package handler

import (
	"encoding/json"
	"net/http"
	"order-service/internal/util/liberr"
	"order-service/internal/util/librest"
	"order-service/module/order/entity"
)

type OrderHandlerConfig struct {
	AuthServiceJWTSecret string
}

type OrderHandler struct {
	orderUsecase OrderUsecase
	configs      OrderHandlerConfig
}

func NewOrderHandler(orderUsecase OrderUsecase, configs OrderHandlerConfig) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
		configs:      configs,
	}
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	user, err := UserAuth(r, o.configs.AuthServiceJWTSecret)
	if err != nil {
		return err
	}

	params := new(entity.CreateOrderRequest)
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		return liberr.NewBaseError(entity.ErrorInvalidBodyJSON)
	}
	params.User = user

	err = o.orderUsecase.CreateOrder(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusCreated
	librest.WriteHTTPResponse(w, entity.GetMessageResponse{
		Message: "Success create order",
		Meta: &entity.Meta{
			HttpStatusCode: code,
		},
	}, code)
	return nil
}
