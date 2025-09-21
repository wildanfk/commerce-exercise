package librest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"warehouse-service/internal/util/librest"

	"github.com/stretchr/testify/assert"
)

func TestGatewayHandlerFunc(t *testing.T) {
	type input struct {
		w  *httptest.ResponseRecorder
		r  *http.Request
		rw *librest.ResponseWriter
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
				handler := librest.GatewayHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
					w.WriteHeader(http.StatusOK)
					return nil
				})

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
				handler := librest.GatewayHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
					w.WriteHeader(http.StatusBadRequest)
					return errors.New("error")
				})

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(i *input) {
				assert.Equal(t, http.StatusBadRequest, i.w.Code)
			},
		},
		{
			name: "Success Handler using Custom Response Writter",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
					w.WriteHeader(http.StatusOK)
					return nil
				})

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				i.rw = librest.NewResponseWriter(i.w)

				handler.ServeHTTP(i.rw, i.r)
			},
			assertFn: func(i *input) {
				assert.Equal(t, http.StatusOK, i.w.Code)
				assert.Equal(t, http.StatusOK, i.rw.StatusCode())
				assert.Nil(t, i.rw.GetError())
			},
		},
		{
			name: "Error Handler using Custom Response Writter",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
					w.WriteHeader(http.StatusBadRequest)
					return errors.New("error")
				})

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				i.rw = librest.NewResponseWriter(i.w)

				handler.ServeHTTP(i.rw, i.r)
			},
			assertFn: func(i *input) {
				assert.Equal(t, http.StatusBadRequest, i.w.Code)
				assert.Equal(t, http.StatusBadRequest, i.rw.StatusCode())
				assert.NotNil(t, i.rw.GetError())
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

func TestApplyGatewayMiddlewares(t *testing.T) {
	middleware1 := func() librest.GatewayMiddleware {
		return func(handle librest.GatewayHandler) librest.GatewayHandler {
			return func(w http.ResponseWriter, r *http.Request) error {
				return handle(w, r)
			}
		}
	}

	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	testCases := []struct {
		name         string
		in           input
		buildInputFn func(*input)
		assertFn     func(*httptest.ResponseRecorder)
	}{
		{
			name: "Success Handler",
			buildInputFn: func(i *input) {
				handler := librest.GatewayHandlerFunc(
					librest.ApplyGatewayMiddlewares(
						func(w http.ResponseWriter, r *http.Request) error {
							w.WriteHeader(http.StatusOK)
							return nil
						}, middleware1(),
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
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
						}, middleware1(),
					),
				)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/handlers", nil)

				handler.ServeHTTP(i.w, i.r)
			},
			assertFn: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildInputFn(&tc.in)
			tc.assertFn(tc.in.w)
		})
	}
}
