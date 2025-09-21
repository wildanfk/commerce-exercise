package usecase_test

import (
	"context"
	"errors"
	"testing"
	"user-service/internal/util/liberr"
	"user-service/module/auth/entity"
	"user-service/module/auth/internal/usecase"
	"user-service/module/auth/internal/usecase/mock"
	"user-service/module/auth/testutil/fixtures"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type authUseCaseDependency struct {
	userRepository *mock.MockUserRepository
	config         *usecase.AuthUsecaseConfig
}

func NewTestAuthUsecase(ctrl *gomock.Controller) (*usecase.AuthUsecase, authUseCaseDependency) {
	useCaseDependency := authUseCaseDependency{
		userRepository: mock.NewMockUserRepository(ctrl),
		config: &usecase.AuthUsecaseConfig{
			JWTSecret:         "secret",
			JWTHourExpiration: 1,
		},
	}

	return usecase.NewAuthUsecase(
		&usecase.AuthUsecaseRepos{
			UserRepo: useCaseDependency.userRepository,
		},
		useCaseDependency.config,
	), useCaseDependency
}

func TestAuth_Authentication(t *testing.T) {
	type input struct {
		params *entity.AuthenticationUserRequest
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*authUseCaseDependency, input)
		assertFn       func(*entity.AuthenticationUser, error)
	}{
		{
			name: "Success Authentication",
			in: input{
				params: &entity.AuthenticationUserRequest{
					Username: "jhon.doe@test.com",
					Password: "test1234",
				},
			},
			mockDependency: func(dependency *authUseCaseDependency, in input) {
				dependency.userRepository.EXPECT().
					GetByUsername(gomock.Any(), in.params.Username).
					Return(fixtures.User, nil)
			},
			assertFn: func(result *entity.AuthenticationUser, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, result)
			},
		},
		{
			name: "Error On Authentication With Wrong Password",
			in: input{
				params: &entity.AuthenticationUserRequest{
					Username: "jhon.doe@test.com",
					Password: "wrongpassword",
				},
			},
			mockDependency: func(dependency *authUseCaseDependency, in input) {
				dependency.userRepository.EXPECT().
					GetByUsername(gomock.Any(), in.params.Username).
					Return(fixtures.User, nil)
			},
			assertFn: func(result *entity.AuthenticationUser, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "Error On GetByUsername",
			in: input{
				params: &entity.AuthenticationUserRequest{
					Username: "jhon.doe@test.com",
					Password: "wrongpassword",
				},
			},
			mockDependency: func(dependency *authUseCaseDependency, in input) {
				dependency.userRepository.EXPECT().
					GetByUsername(gomock.Any(), in.params.Username).
					Return(nil, errors.New("error"))
			},
			assertFn: func(result *entity.AuthenticationUser, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "Error On Authetication With Empty Params",
			in: input{
				params: &entity.AuthenticationUserRequest{},
			},
			mockDependency: func(dependency *authUseCaseDependency, in input) {},
			assertFn: func(result *entity.AuthenticationUser, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)

				berr, ok := err.(*liberr.BaseError)
				assert.True(t, ok)

				details := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("username is a required field", entity.ErrorCodeAuthInvalid, "username"),
					liberr.NewErrorDetails("password is a required field", entity.ErrorCodeAuthInvalid, "password"),
				}
				assert.ElementsMatch(t, details, berr.GetDetails())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()

			ctrl := gomock.NewController(t)
			uc, ucDependency := NewTestAuthUsecase(ctrl)
			defer ctrl.Finish()

			tc.mockDependency(&ucDependency, tc.in)
			tc.assertFn(uc.Authentication(ctx, tc.in.params))
		})
	}
}
