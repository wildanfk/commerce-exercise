package usecase_test

import (
	"context"
	"errors"
	"product-service/module/product/entity"
	"product-service/module/product/internal/usecase"
	"product-service/module/product/internal/usecase/mock"
	"product-service/module/product/testutil/fixtures"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"product-service/internal/util/libpagination"
)

type productUseCaseDependency struct {
	productRepository *mock.MockProductRepository
}

func NewTestProductUsecase(ctrl *gomock.Controller) (*usecase.ProductUsecase, productUseCaseDependency) {
	useCaseDependency := productUseCaseDependency{
		productRepository: mock.NewMockProductRepository(ctrl),
	}

	return usecase.NewProductUsecase(&usecase.ProductUsecaseRepos{
		ProductRepo: useCaseDependency.productRepository,
	}), useCaseDependency
}

func TestProduct_CheckProduct(t *testing.T) {
	type input struct {
		params *entity.ListProductByParams
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*productUseCaseDependency, input)
		assertFn       func([]*entity.Product, *libpagination.OffsetPagination, error)
	}{
		{
			name: "Success Retrieve CheckProduct",
			in: input{
				params: &entity.ListProductByParams{
					Page:  1,
					Limit: 10,
					IDs:   []string{"1", "2"},
				},
			},
			mockDependency: func(dependency *productUseCaseDependency, in input) {
				in.params.Offset = libpagination.Offset(in.params.Page, in.params.Limit)

				dependency.productRepository.EXPECT().
					ListByParams(gomock.Any(), in.params).
					Return([]*entity.Product{fixtures.NewProduct(fixtures.Product)}, &libpagination.OffsetPagination{Offset: 0, Limit: 10, Total: 100}, nil)
			},
			assertFn: func(result []*entity.Product, resultPage *libpagination.OffsetPagination, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, resultPage)
			},
		},
		{
			name: "Error On Retrieve CheckProduct",
			in: input{
				params: &entity.ListProductByParams{
					Page:  1,
					Limit: 10,
					IDs:   []string{"1", "2"},
				},
			},
			mockDependency: func(dependency *productUseCaseDependency, in input) {
				dependency.productRepository.EXPECT().
					ListByParams(gomock.Any(), gomock.Any()).
					Return(nil, nil, errors.New("error"))
			},
			assertFn: func(result []*entity.Product, resultPage *libpagination.OffsetPagination, err error) {
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
			uc, ucDependency := NewTestProductUsecase(ctrl)
			defer ctrl.Finish()

			tc.mockDependency(&ucDependency, tc.in)
			tc.assertFn(uc.CheckProduct(ctx, tc.in.params))
		})
	}
}
