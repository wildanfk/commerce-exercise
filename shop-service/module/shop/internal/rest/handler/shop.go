package handler

import (
	"net/http"
	"shop-service/internal/util"
	"shop-service/internal/util/libpagination"
	"shop-service/internal/util/librest"
	"shop-service/module/shop/entity"
)

const (
	MinimalPageNum  = 1
	MinimalPageSize = 1

	DefaultValueUserListPageNum  = 1
	DefaultValueUserListPageSize = 10
)

type ShopHandler struct {
	shopUsecase ShopUsecase
}

func NewShopHandler(shopUsecase ShopUsecase) *ShopHandler {
	return &ShopHandler{
		shopUsecase: shopUsecase,
	}
}

func (s *ShopHandler) ListByParams(w http.ResponseWriter, r *http.Request) error {
	// Query parameters
	qparams := r.URL.Query()

	params := &entity.ListShopByParams{
		IDs:   qparams["ids"],
		Page:  util.ConvertStringToIntWithDefault(qparams.Get("page_num"), DefaultValueUserListPageNum),
		Limit: util.ConvertStringToIntWithDefault(qparams.Get("page_size"), DefaultValueUserListPageSize),
	}
	if params.Page < MinimalPageNum {
		params.Page = DefaultValueUserListPageNum
	}
	if params.Limit < MinimalPageSize {
		params.Limit = DefaultValueUserListPageSize
	}

	params.Offset = libpagination.Offset(params.Page, params.Limit)

	shops, pagination, err := s.shopUsecase.ListByParams(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.ListShopResponse{
		Shops: shops,
		Meta: &entity.ListMeta{
			Meta: &entity.Meta{
				HttpStatusCode: code,
			},
			PageNum:   pagination.PageNum(),
			PageSize:  pagination.PageSize(),
			PageTotal: pagination.PageTotal(),
		},
	}, code)
	return nil
}
