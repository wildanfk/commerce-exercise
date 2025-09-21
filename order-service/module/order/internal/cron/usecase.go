package cron

import (
	"context"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type OrderUsecase interface {
	ExecuteExpiredOrder(ctx context.Context) error
}
