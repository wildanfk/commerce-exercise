package libcron

import (
	"context"
	"order-service/internal/util/liberr"
	"os"
	"time"

	"go.uber.org/zap"
)

// CronHandler interface that need to be implemented to execute cron
type CronHandler interface {
	ExecuteFunction(ctx context.Context, args []string) error
}

// HandlerFunc helper type to execute cron without implementing interface
type HandlerFunc func(ctx context.Context, args []string) error

// ExecuteFunction implement CronHandler interface
func (f HandlerFunc) ExecuteFunction(ctx context.Context, args []string) error {
	return f(ctx, args)
}

type Config struct {
	Name        string
	CronHandler CronHandler
	Logger      *zap.Logger
}

type Cron struct {
	name        string
	cronHandler CronHandler
	logger      *zap.Logger
}

func NewCron(cfg Config) *Cron {
	cr := &Cron{
		name:   cfg.Name,
		logger: cfg.Logger,
	}

	ch := cfg.CronHandler
	if cr.logger != nil {
		ch = HandlerFunc(cr.executeFunctionWithLogger(ch))
	}

	cr.cronHandler = ch

	return cr
}

func (c *Cron) ExecuteCron() error {
	args := os.Args
	ctx := context.Background()
	return c.cronHandler.ExecuteFunction(ctx, args)
}

func (c *Cron) executeFunctionWithLogger(ch CronHandler) HandlerFunc {
	return func(ctx context.Context, args []string) error {
		timeStart := time.Now()
		err := ch.ExecuteFunction(ctx, args)
		elapsedTime := time.Since(timeStart).Milliseconds()

		fields := []zap.Field{
			zap.String("name", c.name),
			zap.Strings("args", args),
			zap.Int("duration", int(elapsedTime)),
		}

		if err != nil {
			fields = liberr.AppendErrorLogField(fields, err)
			c.logger.Error("Failed executing cron job", fields...)
			return err
		}

		c.logger.Info("Finish executing cron job", fields...)
		return nil
	}
}
