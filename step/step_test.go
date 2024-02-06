package step

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_retry(t *testing.T) {
	t.Run("no error returned - no retry", func(t *testing.T) {
		calledCount := 0
		err := retry(3, func(retriesLeft int) error {
			calledCount++
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 1, calledCount)
	})

	t.Run("single retry, success after", func(t *testing.T) {
		calledCount := 0
		err := retry(3, func(retriesLeft int) error {
			calledCount++
			if calledCount == 1 {
				return errors.New("test error")
			}
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 2, calledCount)
	})

	t.Run("3 attempts, all fails - error returned", func(t *testing.T) {
		calledCount := 0
		err := retry(3, func(retriesLeft int) error {
			calledCount++
			return errors.New("all attempts failed")
		})
		require.EqualError(t, err, "all attempts failed")
		assert.Equal(t, 3, calledCount)
	})
}
