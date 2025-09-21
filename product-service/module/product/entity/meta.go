package entity

type Meta struct {
	HttpStatusCode int `json:"http_status_code"`
}

type ListMeta struct {
	*Meta
	PageNum   int `json:"page_num"`
	PageSize  int `json:"page_size"`
	PageTotal int `json:"page_total"`
}
