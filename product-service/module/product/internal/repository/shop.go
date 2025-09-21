package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"product-service/internal/util/liberr"
	"product-service/module/product/entity"
	"strconv"
	"time"

	rh "github.com/hashicorp/go-retryablehttp"
)

type ShopConfiguration struct {
	ApiHost           string
	BasicAuthUsername string
	BasicAuthPassword string
}

type ShopRepository struct {
	Config     ShopConfiguration
	httpClient *rh.Client
}

func NewShopRepository(config ShopConfiguration, httpClient *rh.Client) *ShopRepository {
	return &ShopRepository{
		Config:     config,
		httpClient: httpClient,
	}
}

type shop struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type listShopResponse struct {
	Shops []*shop          `json:"shops"`
	Meta  *entity.ListMeta `json:"meta"`
}

func (s *ShopRepository) basicAuth() string {
	auth := s.Config.BasicAuthUsername + ":" + s.Config.BasicAuthPassword
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encodedAuth
}

func (w *ShopRepository) ListByShopIDs(ctx context.Context, shopIDs []string) ([]*entity.Shop, error) {
	qparams := url.Values{}
	for _, sid := range shopIDs {
		qparams.Add("ids", sid)
	}
	qparams.Add("page_size", strconv.Itoa(len(shopIDs)))

	path := w.Config.ApiHost + "/shops?" + qparams.Encode()

	req, _ := rh.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", w.basicAuth())

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, liberr.NewTracer("Error when request on Shop.ListByParams").Wrap(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, liberr.NewTracer("Error happened when read body response on Shop.ListByParams").Wrap(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, liberr.NewTracer(fmt.Sprintf("Error with http status %d on Shop.ListByParams", resp.StatusCode)).Wrap(err)
	}

	responseObj := listShopResponse{}
	err = json.Unmarshal(responseBody, &responseObj)
	if err != nil {
		return nil, liberr.NewTracer("Error happened when parse json body response on Shop.ListByParams").Wrap(err)
	}

	shops := []*entity.Shop{}
	if responseObj.Shops != nil {
		for _, s := range responseObj.Shops {
			shops = append(shops, &entity.Shop{
				ID:   s.ID,
				Name: s.Name,
			})
		}
	}

	return shops, nil
}
