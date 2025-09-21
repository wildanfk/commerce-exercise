package entity

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListProductRequest struct {
	PageNum  int
	PageSize int
	IDs      []string
}

type ListProductResponse struct {
	Products []*Product `json:"products"`
	Meta     *ListMeta  `json:"meta"`
}
