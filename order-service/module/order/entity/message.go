package entity

type GetMessageResponse struct {
	Message string `json:"message"`
	Meta    *Meta  `json:"meta"`
}
