package util_test

import (
	"testing"
	"warehouse-service/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestNowUTCWithoutNanoSecond(t *testing.T) {
	tests := []struct {
		name     string
		assertFn func()
	}{
		{
			name: "Success Retrieve NowUTCWithoutNanoSecond",
			assertFn: func() {
				assert.NotNil(t, util.NowUTCWithoutNanoSecond())
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn()
		})
	}
}
