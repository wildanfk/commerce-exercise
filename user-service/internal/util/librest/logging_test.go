package librest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/util/librest"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithLoggingMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	testCases := []struct {
		name         string
		in           input
		buildInputFn func(*input)
		assertFn     func(*input)
	}{
		{
			name: "Success Handler",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							w.WriteHeader(http.StatusOK)
							return nil
						}, librest.WithLoggingMiddleware("/handlers", logger),
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				assert.Equal(t, http.StatusOK, i.w.Code)
			},
		},
		{
			name: "Error Handler",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							w.WriteHeader(http.StatusBadRequest)
							return errors.New("error")
						}, librest.WithLoggingMiddleware("/handlers", logger),
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				assert.Equal(t, http.StatusBadRequest, i.w.Code)
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
