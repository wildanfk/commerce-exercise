package util_test

import (
	"testing"
	"user-service/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestHashPlainTextWithSHA256(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		assertFn  func(string)
	}{
		{
			name:      "Success HashPlainTextWithSHA256 with value test1234",
			plaintext: "test1234",
			assertFn: func(result string) {
				assert.Equal(t, "937e8d5fbb48bd4949536cd65b8d35c426b80d2f830c5c308e2cdec422ae2244", result)
			},
		},
		{
			name:      "Success HashPlainTextWithSHA256 with value jhondoe12345loremipsum",
			plaintext: "jhondoe12345loremipsum",
			assertFn: func(result string) {
				assert.Equal(t, "b0868ad5f548443eee872d9a2a0e04d6e304adb0ca437e1e008ad5d5a5ed61f3", result)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(util.HashPlainTextWithSHA256(tc.plaintext))
		})
	}
}
