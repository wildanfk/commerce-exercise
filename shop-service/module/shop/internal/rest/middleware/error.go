package middleware

import (
	"net/http"
	"shop-service/internal/util/liberr"
	"shop-service/internal/util/librest"
	"shop-service/module/shop/entity"
)

var (
	errorCodeMapper = map[string]int{
		entity.ErrorCodeForbidden:    http.StatusForbidden,
		entity.ErrorCodeShopNotFound: http.StatusNotFound,
	}
)

type errorMiddleware struct {
	handler librest.GatewayHandler
}

func WithErrorMiddleware() librest.GatewayMiddleware {
	return func(handle librest.GatewayHandler) librest.GatewayHandler {
		em := errorMiddleware{handler: handle}

		return em.handle
	}
}

func (em *errorMiddleware) bodyFromBaseError(berr *liberr.BaseError) entity.ErrorResponse {
	errors := []*entity.Error{}
	for _, detail := range berr.GetDetails() {
		errors = append(errors, &entity.Error{
			ErrorCode:    detail.Code,
			ErrorMessage: detail.Message,
			ErrorField:   detail.Field,
		})
	}

	status := http.StatusBadRequest
	for errCode, errStatus := range errorCodeMapper {
		if berr.IsAnyCodeEqual(errCode) {
			status = errStatus
			break
		}
	}

	return entity.ErrorResponse{
		Errors: errors,
		Meta: &entity.Meta{
			HttpStatusCode: status,
		},
	}
}

func (em *errorMiddleware) bodyFromAnyError() entity.ErrorResponse {
	return entity.ErrorResponse{
		Errors: []*entity.Error{
			{
				ErrorCode:    "INTERNAL_SERVER_ERROR",
				ErrorMessage: "Internal Server Error",
			},
		},
		Meta: &entity.Meta{
			HttpStatusCode: http.StatusInternalServerError,
		},
	}
}

func (em *errorMiddleware) handle(w http.ResponseWriter, r *http.Request) error {
	rw := librest.WrapResponseWriter(w)
	err := em.handler(rw, r)
	if err != nil {
		var body entity.ErrorResponse

		if berr, ok := err.(*liberr.BaseError); ok {
			body = em.bodyFromBaseError(berr)
		} else {
			body = em.bodyFromAnyError()
		}

		librest.WriteHTTPResponse(rw, body, body.Meta.HttpStatusCode)
	}

	return err
}
