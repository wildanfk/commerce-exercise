package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/testutil"
	"user-service/module/auth/entity"
	"user-service/module/auth/internal/rest/handler"
	"user-service/module/auth/internal/rest/handler/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type authHandlerDependency struct {
	authUsecase *mock.MockAuthUsecase
}

func NewTestAuthHandler(ctrl *gomock.Controller) (*handler.AuthHandler, authHandlerDependency) {
	userHandlerDependency := authHandlerDependency{
		authUsecase: mock.NewMockAuthUsecase(ctrl),
	}

	return handler.NewAuthHandler(userHandlerDependency.authUsecase), userHandlerDependency
}

func TestAuthHandler_Authentication(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	testCases := []struct {
		name           string
		in             input
		buildInputFn   func(*input)
		mockDependency func(*authHandlerDependency, input)
		assertFn       func(*httptest.ResponseRecorder, error)
	}{
		{
			name: "Success Authentication",
			buildInputFn: func(i *input) {
				params := entity.AuthenticationUserRequest{
					Username: "jhon.doe@test.com",
					Password: "test1234",
				}
				body, _ := json.Marshal(params)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodPost, "/authentication", bytes.NewBuffer(body))
			},
			mockDependency: func(dependency *authHandlerDependency, i input) {
				dependency.authUsecase.EXPECT().
					Authentication(
						gomock.Any(),
						&entity.AuthenticationUserRequest{
							Username: "jhon.doe@test.com",
							Password: "test1234",
						}).
					Return(&entity.AuthenticationUser{
						AccessToken: "token",
					}, nil)
			},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.Nil(t, err)

				expected := entity.AuthenticationUserResponse{
					Auth: &entity.AuthenticationUser{
						AccessToken: "token",
					},
					Meta: &entity.Meta{
						HttpStatusCode: http.StatusOK,
					},
				}

				assert.Equal(t, expected.Meta.HttpStatusCode, rr.Code)

				var actual entity.AuthenticationUserResponse
				_ = json.NewDecoder(rr.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Error On Authentication",
			buildInputFn: func(i *input) {
				params := entity.AuthenticationUserRequest{
					Username: "jhon.doe@test.com",
					Password: "test1234",
				}
				body, _ := json.Marshal(params)

				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodPost, "/authentication", bytes.NewBuffer(body))
			},
			mockDependency: func(dependency *authHandlerDependency, i input) {
				dependency.authUsecase.EXPECT().
					Authentication(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error On Body Request",
			buildInputFn: func(i *input) {
				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodPost, "/authentication", testutil.ErrReader(0))
			},
			mockDependency: func(dependency *authHandlerDependency, i input) {},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			handler, handlerDependency := NewTestAuthHandler(ctrl)
			defer ctrl.Finish()

			tc.buildInputFn(&tc.in)
			tc.mockDependency(&handlerDependency, tc.in)

			err := handler.Authentication(tc.in.w, tc.in.r)
			tc.assertFn(tc.in.w, err)
		})
	}
}
