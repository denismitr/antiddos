package challenge_test

import (
	"github.com/denismitr/antiddos/internal/challenge"
	"github.com/denismitr/antiddos/internal/store/adapters/nope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChallenge_Solve(t *testing.T) {
	t.Run("3 zeroes", func(t *testing.T) {
		c := challenge.New(nope.Nope{}, 3, 30)
		c.SetNow(func() time.Time {
			return time.Unix(1702740115, 0)
		})
		header, err := c.Solve("1|3|1702740115|127.0.0.1:52374|ODk1Mw==|0")
		require.NoError(t, err)
		assert.Equal(t, "1|3|1702740115|127.0.0.1:52374|ODk1Mw==|2797", header)
	})

	t.Run("3 zeroes", func(t *testing.T) {
		c := challenge.New(nope.Nope{}, 3, 30)
		require.NotNil(t, c)
		c.SetNow(func() time.Time {
			return time.Unix(1702740115, 0)
		})
		c.SetRandomizer(func() int {
			return 5000
		})

		header, err := c.Create("hello world!")
		require.NoError(t, err)
		assert.Equal(t, "1|3|1702740115|hello world!|NTAwMA==|0", header)
	})
}
