package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"shop-service/internal/util/liberr"
	"shop-service/internal/util/librest"
	"shop-service/module/shop/entity"
	"shop-service/module/shop/internal/rest/middleware"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorMiddleware(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	middlewares := []librest.GatewayMiddleware{
		middleware.WithErrorMiddleware(),
	}

	testCases := []struct {
		name         string
		in           input
		buildInputFn func(*input)
		assertFn     func(*input)
	}{
		{
			name: "Error Handler With Base Error",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							return liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "field"))
						}, middlewares...,
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				expected := &entity.ErrorResponse{
					Errors: []*entity.Error{
						{ErrorMessage: "Error Happened", ErrorCode: "ERROR_CODE", ErrorField: "field"},
					},
					Meta: &entity.Meta{
						HttpStatusCode: http.StatusBadRequest,
					},
				}

				assert.Equal(t, expected.Meta.HttpStatusCode, i.w.Code)

				var actual *entity.ErrorResponse
				_ = json.NewDecoder(i.w.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Error Handler With Base Error That Has Code On errorCodeMapper",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							return liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", entity.ErrorCodeShopNotFound, "field"))
						}, middlewares...,
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				expected := &entity.ErrorResponse{
					Errors: []*entity.Error{
						{ErrorMessage: "Error Happened", ErrorCode: entity.ErrorCodeShopNotFound, ErrorField: "field"},
					},
					Meta: &entity.Meta{
						HttpStatusCode: http.StatusNotFound,
					},
				}

				assert.Equal(t, expected.Meta.HttpStatusCode, i.w.Code)

				var actual *entity.ErrorResponse
				_ = json.NewDecoder(i.w.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Error Handler With Any Error",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							return errors.New("Error Happened")
						}, middlewares...,
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				expected := &entity.ErrorResponse{
					Errors: []*entity.Error{
						{ErrorMessage: "Internal Server Error", ErrorCode: "INTERNAL_SERVER_ERROR", ErrorField: ""},
					},
					Meta: &entity.Meta{
						HttpStatusCode: http.StatusInternalServerError,
					},
				}

				assert.Equal(t, expected.Meta.HttpStatusCode, i.w.Code)

				var actual *entity.ErrorResponse
				_ = json.NewDecoder(i.w.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildInputFn(&tc.in)
			tc.assertFn(&tc.in)
		})
	}
}
