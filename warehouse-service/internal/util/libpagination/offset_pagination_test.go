package libpagination_test

import (
	"testing"
	"warehouse-service/internal/util/libpagination"

	"github.com/stretchr/testify/assert"
)

func TestOffsetPagination_PageSize(t *testing.T) {
	type input struct {
		op libpagination.OffsetPagination
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(int)
	}{
		{
			name: "Success Calculate with Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Limit: 10,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
		{
			name: "Success Calculate with Limit 20",
			in: input{
				op: libpagination.OffsetPagination{
					Limit: 20,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 20, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.op.PageSize())
		})
	}
}

func TestOffsetPagination_PageNum(t *testing.T) {
	type input struct {
		op libpagination.OffsetPagination
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(int)
	}{
		{
			name: "Success Calculate with Offset 0 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 0,
					Limit:  10,
					Total:  100,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 1, result)
			},
		},
		{
			name: "Success Calculate with Offset 10 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 9,
					Limit:  10,
					Total:  100,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 1, result)
			},
		},
		{
			name: "Success Calculate with Offset 11 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 10,
					Limit:  10,
					Total:  100,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 2, result)
			},
		},
		{
			name: "Success Calculate with Offset 97 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 97,
					Limit:  10,
					Total:  100,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.op.PageNum())
		})
	}
}

func TestOffsetPagination_PageTotal(t *testing.T) {
	type input struct {
		op libpagination.OffsetPagination
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(int)
	}{
		{
			name: "Success Calculate with Total 100 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 0,
					Limit:  10,
					Total:  100,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
		{
			name: "Success Calculate with Total 103 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 0,
					Limit:  10,
					Total:  103,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 11, result)
			},
		},
		{
			name: "Success Calculate with Total 10 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 0,
					Limit:  10,
					Total:  10,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 1, result)
			},
		},
		{
			name: "Success Calculate with Total 5 & Limit 10",
			in: input{
				op: libpagination.OffsetPagination{
					Offset: 0,
					Limit:  10,
					Total:  5,
				},
			},
			assertFn: func(result int) {
				assert.Equal(t, 1, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.op.PageTotal())
		})
	}
}

func TestOffset(t *testing.T) {
	type input struct {
		pagenum int
		limit   int
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(int)
	}{
		{
			name: "Success Calculate with Page 1 & Limit 10",
			in: input{
				pagenum: 1,
				limit:   10,
			},
			assertFn: func(result int) {
				assert.Equal(t, 0, result)
			},
		},
		{
			name: "Success Calculate with Page 2 & Limit 10",
			in: input{
				pagenum: 2,
				limit:   10,
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
		{
			name: "Success Calculate with Page 10 & Limit 10",
			in: input{
				pagenum: 10,
				limit:   10,
			},
			assertFn: func(result int) {
				assert.Equal(t, 90, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(libpagination.Offset(tc.in.pagenum, tc.in.limit))
		})
	}
}
