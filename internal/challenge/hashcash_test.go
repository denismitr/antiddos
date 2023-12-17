package challenge

import (
	"encoding/base64"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashcash_Bruteforce(t *testing.T) {
	now := time.Date(2023, 12, 15, 10, 45, 20, 0, time.UTC)

	t.Run("4 zeros", func(t *testing.T) {
		data := hashcash{
			Ver:      1,
			Bits:     4,
			Date:     uint64(now.Unix()),
			Resource: "some transmitted data",
			Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 467124))),
			Counter:  0,
		}
		err := data.Bruteforce(math.MaxInt)
		require.NoError(t, err)
		assert.Equal(t, 8879, int(data.Counter))
	})

	t.Run("5 zeros", func(t *testing.T) {
		data := hashcash{
			Ver:      1,
			Bits:     5,
			Date:     uint64(now.Unix()),
			Resource: "some transmitted data",
			Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 557399))),
			Counter:  0,
		}
		err := data.Bruteforce(math.MaxUint64)
		require.NoError(t, err)
		assert.Equal(t, 1037588, int(data.Counter))
	})

	t.Run("too many zeroes", func(t *testing.T) {
		hc := hashcash{
			Ver:      1,
			Bits:     10,
			Date:     uint64(now.Unix()),
			Resource: "some transmitted data",
			Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:  0,
		}

		const iterations = 300_000
		err := hc.Bruteforce(iterations)
		require.Error(t, err)
		assert.Equal(t, iterations+1, int(hc.Counter))
		//assert.True(t, errors.Is(ErrTooManyIterations, err))
	})
}
