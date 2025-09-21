package libpagination

import "math"

type OffsetPagination struct {
	Offset int
	Limit  int
	Total  int
}

func (op OffsetPagination) PageSize() int {
	return op.Limit
}

func (op OffsetPagination) PageNum() int {
	return (op.Offset / op.Limit) + 1
}

func (op OffsetPagination) PageTotal() int {
	return int(math.Ceil(float64(op.Total) / float64(op.Limit)))
}

func Offset(pagenum int, limit int) int {
	return (pagenum - 1) * limit
}
