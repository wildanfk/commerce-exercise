package librest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-service/internal/util/librest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteHTTPResponse(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}

	type input struct {
		body any
		code int
	}

	tests := []struct {
		name     string
		in       input
		assertFn func(*httptest.ResponseRecorder, input)
	}{
		{
			name: "Success On Response OK",
			in: input{
				body: response{Message: "OK"},
				code: 200,
			},
			assertFn: func(rr *httptest.ResponseRecorder, in input) {
				assert.Equal(t, in.code, rr.Code)

				var actual response
				_ = json.NewDecoder(rr.Body).Decode(&actual)
				assert.Equal(t, in.body, actual)
			},
		},
		{
			name: "Success On Bad Response",
			in: input{
				body: response{Message: "Bad Request"},
				code: 400,
			},
			assertFn: func(rr *httptest.ResponseRecorder, in input) {
				assert.Equal(t, in.code, rr.Code)

				var actual response
				_ = json.NewDecoder(rr.Body).Decode(&actual)
				assert.Equal(t, in.body, actual)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			librest.WriteHTTPResponse(rr, tc.in.body, tc.in.code)

			tc.assertFn(rr, tc.in)
		})
	}
}
