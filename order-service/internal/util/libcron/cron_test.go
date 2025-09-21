package libcron_test

import (
	"context"
	"errors"
	"order-service/internal/util/libcron"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type cronHandler struct {
	err error
}

func (c cronHandler) ExecuteFunction(ctx context.Context, args []string) error {
	return c.err
}

func TestCron_ExecuteCron(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type input struct {
		cron *libcron.Cron
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(error)
	}{
		{
			name: "Success on ExecuteCron",
			in: input{
				cron: func() *libcron.Cron {
					ch := cronHandler{}

					return libcron.NewCron(
						libcron.Config{
							Name:        "sample_cron",
							Logger:      logger,
							CronHandler: ch,
						},
					)
				}(),
			},
			assertFn: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "Error on ExecuteCron",
			in: input{
				cron: func() *libcron.Cron {
					ch := cronHandler{
						err: errors.New("error happened"),
					}

					return libcron.NewCron(
						libcron.Config{
							Name:        "sample_cron",
							Logger:      logger,
							CronHandler: ch,
						},
					)
				}(),
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.cron.ExecuteCron())
		})
	}
}
