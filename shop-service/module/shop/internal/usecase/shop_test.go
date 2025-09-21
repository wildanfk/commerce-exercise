package usecase_test

import (
	"context"
	"errors"
	"shop-service/module/shop/entity"
	"shop-service/module/shop/internal/usecase"
	"shop-service/module/shop/internal/usecase/mock"
	"shop-service/module/shop/testutil/fixtures"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"shop-service/internal/util/libpagination"
)

type shopUseCaseDependency struct {
	shopRepository *mock.MockShopRepository
}

func NewTestShopUsecase(ctrl *gomock.Controller) (*usecase.ShopUsecase, shopUseCaseDependency) {
	useCaseDependency := shopUseCaseDependency{
		shopRepository: mock.NewMockShopRepository(ctrl),
	}

	return usecase.NewShopUsecase(&usecase.ShopUsecaseRepos{
		ShopRepo: useCaseDependency.shopRepository,
	}), useCaseDependency
}

func TestUser_ListByParams(t *testing.T) {
	type input struct {
		params *entity.ListShopByParams
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*shopUseCaseDependency, input)
		assertFn       func([]*entity.Shop, *libpagination.OffsetPagination, error)
	}{
		{
			name: "Success Retrieve ListByParams",
			in: input{
				params: &entity.ListShopByParams{
					Page:  1,
					Limit: 10,
					IDs:   []string{"1", "2"},
				},
			},
			mockDependency: func(dependency *shopUseCaseDependency, in input) {
				in.params.Offset = libpagination.Offset(in.params.Page, in.params.Limit)

				dependency.shopRepository.EXPECT().
					ListByParams(gomock.Any(), in.params).
					Return([]*entity.Shop{fixtures.NewShop(fixtures.Shop)}, &libpagination.OffsetPagination{Offset: 0, Limit: 10, Total: 100}, nil)
			},
			assertFn: func(result []*entity.Shop, resultPage *libpagination.OffsetPagination, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, resultPage)
			},
		},
		{
			name: "Error On Retrieve ListByParams",
			in: input{
				params: &entity.ListShopByParams{
					Page:  1,
					Limit: 10,
					IDs:   []string{"1", "2"},
				},
			},
			mockDependency: func(dependency *shopUseCaseDependency, in input) {
				dependency.shopRepository.EXPECT().
					ListByParams(gomock.Any(), gomock.Any()).
					Return(nil, nil, errors.New("error"))
			},
			assertFn: func(result []*entity.Shop, resultPage *libpagination.OffsetPagination, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Nil(t, resultPage)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()

			ctrl := gomock.NewController(t)
			uc, ucDependency := NewTestShopUsecase(ctrl)
			defer ctrl.Finish()

			tc.mockDependency(&ucDependency, tc.in)
			tc.assertFn(uc.ListByParams(ctx, tc.in.params))
		})
	}
}
