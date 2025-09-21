package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"shop-service/internal/util/libpagination"
	"shop-service/module/shop/entity"
	"shop-service/module/shop/internal/rest/handler"
	"shop-service/module/shop/internal/rest/handler/mock"
	"shop-service/module/shop/testutil/fixtures"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type shopHandlerDependency struct {
	shopUsecase *mock.MockShopUsecase
}

func NewTestShopHandler(ctrl *gomock.Controller) (*handler.ShopHandler, shopHandlerDependency) {
	shopHandlerDependency := shopHandlerDependency{
		shopUsecase: mock.NewMockShopUsecase(ctrl),
	}

	return handler.NewShopHandler(shopHandlerDependency.shopUsecase), shopHandlerDependency
}

func TestShopHandler_ListByParams(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	testCases := []struct {
		name           string
		in             input
		buildInputFn   func(*input)
		mockDependency func(*shopHandlerDependency, input)
		assertFn       func(*httptest.ResponseRecorder, error)
	}{
		{
			name: "Success Retrieve ListByParams",
			buildInputFn: func(i *input) {
				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/shops?page_num=3&page_size=10&ids=2&ids=4", nil)
			},
			mockDependency: func(dependency *shopHandlerDependency, i input) {
				dependency.shopUsecase.EXPECT().
					ListByParams(gomock.Any(), &entity.ListShopByParams{
						Page:   3,
						Limit:  10,
						Offset: 20,
						IDs:    []string{"2", "4"},
					}).
					Return(
						[]*entity.Shop{fixtures.NewShop(fixtures.Shop)},
						&libpagination.OffsetPagination{
							Offset: 20,
							Limit:  10,
							Total:  100,
						},
						nil)
			},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, rr.Code)

				expected := entity.ListShopResponse{
					Shops: []*entity.Shop{fixtures.NewShop(fixtures.Shop)},
					Meta: &entity.ListMeta{
						Meta: &entity.Meta{
							HttpStatusCode: http.StatusOK,
						},
						PageNum:   3,
						PageSize:  10,
						PageTotal: 10,
					},
				}

				var actual entity.ListShopResponse
				_ = json.NewDecoder(rr.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Success Retrieve ListByParams With Empty Params",
			buildInputFn: func(i *input) {
				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/shops", nil)
			},
			mockDependency: func(dependency *shopHandlerDependency, i input) {
				dependency.shopUsecase.EXPECT().
					ListByParams(gomock.Any(), &entity.ListShopByParams{
						Page:   1,
						Limit:  10,
						Offset: 0,
					}).
					Return(
						[]*entity.Shop{fixtures.NewShop(fixtures.Shop)},
						&libpagination.OffsetPagination{
							Offset: 0,
							Limit:  10,
							Total:  100,
						},
						nil)
			},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, rr.Code)

				expected := entity.ListShopResponse{
					Shops: []*entity.Shop{fixtures.NewShop(fixtures.Shop)},
					Meta: &entity.ListMeta{
						Meta: &entity.Meta{
							HttpStatusCode: http.StatusOK,
						},
						PageNum:   1,
						PageSize:  10,
						PageTotal: 10,
					},
				}

				var actual entity.ListShopResponse
				_ = json.NewDecoder(rr.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Success Retrieve ListByParams With Invalid Page Params",
			buildInputFn: func(i *input) {
				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/shops?page_num=0&page_size=0", nil)
			},
			mockDependency: func(dependency *shopHandlerDependency, i input) {
				dependency.shopUsecase.EXPECT().
					ListByParams(gomock.Any(), &entity.ListShopByParams{
						Page:   1,
						Limit:  10,
						Offset: 0,
					}).
					Return(
						[]*entity.Shop{fixtures.NewShop(fixtures.Shop)},
						&libpagination.OffsetPagination{
							Offset: 0,
							Limit:  10,
							Total:  100,
						},
						nil)
			},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, rr.Code)

				expected := entity.ListShopResponse{
					Shops: []*entity.Shop{fixtures.NewShop(fixtures.Shop)},
					Meta: &entity.ListMeta{
						Meta: &entity.Meta{
							HttpStatusCode: http.StatusOK,
						},
						PageNum:   1,
						PageSize:  10,
						PageTotal: 10,
					},
				}

				var actual entity.ListShopResponse
				_ = json.NewDecoder(rr.Body).Decode(&actual)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Error On Retrieve ListByParams",
			buildInputFn: func(i *input) {
				i.w = httptest.NewRecorder()
				i.r = httptest.NewRequest(http.MethodGet, "/shops", nil)
			},
			mockDependency: func(dependency *shopHandlerDependency, i input) {
				dependency.shopUsecase.EXPECT().
					ListByParams(gomock.Any(), gomock.Any()).
					Return(nil, nil, errors.New("error"))
			},
			assertFn: func(rr *httptest.ResponseRecorder, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			handler, handlerDependency := NewTestShopHandler(ctrl)
			defer ctrl.Finish()

			tc.buildInputFn(&tc.in)
			tc.mockDependency(&handlerDependency, tc.in)

			err := handler.ListByParams(tc.in.w, tc.in.r)
			tc.assertFn(tc.in.w, err)
		})
	}
}
