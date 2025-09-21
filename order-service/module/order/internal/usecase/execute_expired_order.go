package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func (o *OrderUsecase) ExecuteExpiredOrder(ctx context.Context) error {
	logFields := []zap.Field{
		zap.String("function", "ExecuteExpiredOrder"),
	}

	expiredOrders, err := o.repos.OrderRepo.ListByOrderExpired(ctx)
	if err != nil {
		return err
	}

	for _, eo := range expiredOrders {
		o.logger.Info(fmt.Sprintf("Expired Order ID : %s", eo.ID), logFields...)

		orderDetails, err := o.repos.OrderDetailRepo.ListByOrderID(ctx, eo.ID)
		if err != nil {
			o.logger.Error(fmt.Sprintf("Expired Order ID : %s Failed on Retrieve Order Detail due %v", eo.ID, err), logFields...)
			continue
		}

		tx, err := o.repos.DatabaseTransactionHandler.Begin(ctx, &sql.TxOptions{})
		if err != nil {
			o.logger.Error(fmt.Sprintf("Expired Order ID : %s Failed on Begin Transaction due %v", eo.ID, err), logFields...)
			continue
		}

		_, err = o.repos.OrderRepo.UpdateExpired(ctx, eo.ID, tx)
		if err != nil {
			tx.Rollback() //nolint
			o.logger.Error(fmt.Sprintf("Expired Order ID : %s Failed on UpdateExpired due %v", eo.ID, err), logFields...)
			continue
		}

		err = o.releaseStocks(ctx, orderDetails)
		if err != nil {
			tx.Rollback() //nolint
			o.logger.Error(fmt.Sprintf("Expired Order ID : %s Failed on ReleaseStock due %v", eo.ID, err), logFields...)
			continue
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback() //nolint
			o.logger.Error(fmt.Sprintf("Expired Order ID : %s Failed on Commit due %v", eo.ID, err), logFields...)
		}
	}

	return nil
}
