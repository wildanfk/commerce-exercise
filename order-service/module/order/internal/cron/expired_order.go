package cron

import (
	"context"
)

type ExpiredOrderCron struct {
	orderUsecase OrderUsecase
}

func NewExpiredOrderCron(orderUsecase OrderUsecase) *ExpiredOrderCron {
	return &ExpiredOrderCron{
		orderUsecase: orderUsecase,
	}
}

func (e ExpiredOrderCron) ExecuteFunction(ctx context.Context, args []string) error {
	return e.orderUsecase.ExecuteExpiredOrder(ctx)
}
