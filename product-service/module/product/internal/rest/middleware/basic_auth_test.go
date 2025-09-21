package middleware_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-service/internal/util/librest"
	"product-service/module/product/entity"
	"product-service/module/product/internal/rest/middleware"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithBasicAuthMiddleware(t *testing.T) {
	basicAuthUser := "auth_user"
	basicAuthPass := "auth_pass"

	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	middlewares := []librest.GatewayMiddleware{
		middleware.WithBasicAuthMiddleware(basicAuthUser, basicAuthPass),
		middleware.WithErrorMiddleware(),
	}

	testCases := []struct {
		name         string
		in           input
		buildInputFn func(*input)
		assertFn     func(*input)
	}{
		{
			name: "Success Basic Auth Handler",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							w.WriteHeader(http.StatusOK)
							return nil
						}, middlewares...,
					),
				)

				auth := basicAuthUser + ":" + basicAuthPass
				encoded := base64.StdEncoding.EncodeToString([]byte(auth))
				headerValue := "Basic " + encoded

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)
				i.r.Header.Set("Authorization", headerValue)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				assert.Equal(t, http.StatusOK, i.w.Code)
			},
		},
		{
			name: "Failed Basic Auth Handler",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							w.WriteHeader(http.StatusOK)
							return nil
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
						{ErrorMessage: "Forbidden", ErrorCode: entity.ErrorCodeForbidden, ErrorField: ""},
					},
					Meta: &entity.Meta{
						HttpStatusCode: http.StatusForbidden,
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
