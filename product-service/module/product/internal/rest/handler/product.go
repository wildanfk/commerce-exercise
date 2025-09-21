package handler

import (
	"net/http"
	"product-service/internal/util"
	"product-service/internal/util/libpagination"
	"product-service/internal/util/librest"
	"product-service/module/product/entity"
)

const (
	MinimalPageNum  = 1
	MinimalPageSize = 1

	DefaultValueUserListPageNum  = 1
	DefaultValueUserListPageSize = 10
)

type ProductHandler struct {
	productUsecase ProductUsecase
}

func NewProductHandler(productUsecase ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

func (p *ProductHandler) CheckProduct(w http.ResponseWriter, r *http.Request) error {
	// Query parameters
	qparams := r.URL.Query()

	params := &entity.ListProductByParams{
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

	products, pagination, err := p.productUsecase.CheckProduct(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.ListProductResponse{
		Products: products,
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

func (p *ProductHandler) ListProduct(w http.ResponseWriter, r *http.Request) error {
	// Query parameters
	qparams := r.URL.Query()

	params := &entity.ListProductByParams{
		Name:  qparams.Get("name"),
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

	products, pagination, err := p.productUsecase.ListProduct(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.ListProductDetailResponse{
		Products: products,
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
