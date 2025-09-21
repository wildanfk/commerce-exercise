package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"order-service/internal/util/liberr"
	"order-service/module/order/entity"
	"strconv"
	"time"

	rh "github.com/hashicorp/go-retryablehttp"
	"github.com/shopspring/decimal"
)

type ProductConfiguration struct {
	ApiHost           string
	BasicAuthUsername string
	BasicAuthPassword string
}

type ProductRepository struct {
	Config     ProductConfiguration
	httpClient *rh.Client
}

func NewProductRepository(config ProductConfiguration, httpClient *rh.Client) *ProductRepository {
	return &ProductRepository{
		Config:     config,
		httpClient: httpClient,
	}
}

type product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type listProductResponse struct {
	Products []*product       `json:"products"`
	Meta     *entity.ListMeta `json:"meta"`
}

func (p *ProductRepository) basicAuth() string {
	auth := p.Config.BasicAuthUsername + ":" + p.Config.BasicAuthPassword
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encodedAuth
}

func (p *ProductRepository) ListByProductIDs(ctx context.Context, productIDs []string) ([]*entity.Product, error) {
	qparams := url.Values{}
	for _, pid := range productIDs {
		qparams.Add("ids", pid)
	}
	qparams.Add("page_size", strconv.Itoa(len(productIDs)))

	path := p.Config.ApiHost + "/check-products?" + qparams.Encode()

	req, _ := rh.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", p.basicAuth())

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, liberr.NewTracer("Error when request on Product.ListByProductIDs").Wrap(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, liberr.NewTracer("Error happened when read body response on Product.ListByProductIDs").Wrap(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, liberr.NewTracer(fmt.Sprintf("Error with http status %d on Product.ListByProductIDs", resp.StatusCode)).Wrap(err)
	}

	responseObj := listProductResponse{}
	err = json.Unmarshal(responseBody, &responseObj)
	if err != nil {
		return nil, liberr.NewTracer("Error happened when parse json body response on Product.ListByProductIDs").Wrap(err)
	}

	products := []*entity.Product{}
	if responseObj.Products != nil {
		for _, p := range responseObj.Products {
			price, _ := decimal.NewFromString(p.Price)

			products = append(products, &entity.Product{
				ID:    p.ID,
				Name:  p.Name,
				Price: price,
			})
		}
	}

	return products, nil
}
