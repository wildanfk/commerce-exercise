package entity

import "time"

type Shop struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListShopByParams struct {
	Page   int
	Offset int
	Limit  int
	IDs    []string
}

type ListShopResponse struct {
	Shops []*Shop   `json:"shops"`
	Meta  *ListMeta `json:"meta"`
}
