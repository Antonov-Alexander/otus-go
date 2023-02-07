package storages

import (
	"testing"
	"time"

	"github.com/Antonov-Alexander/otus-go/final_project/types"
	"github.com/stretchr/testify/require"
)

func TestCalcBucket(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		testCases := []struct {
			input    types.Limit
			expected BucketInfo
		}{
			{
				input: types.Limit{
					Interval: 30,
					Limit:    1,
				},
				expected: BucketInfo{
					Interval: 30000,
					Limit:    1,
					Count:    1,
				},
			},
			{
				input: types.Limit{
					Interval: 30,
					Limit:    15,
				},
				expected: BucketInfo{
					Interval: 2000,
					Limit:    1,
					Count:    15,
				},
			},
			{
				input: types.Limit{
					Interval: 15,
					Limit:    30,
				},
				expected: BucketInfo{
					Interval: 1500,
					Limit:    3,
					Count:    10,
				},
			},
		}

		memoryStorage := MemoryStorage{}
		memoryStorage.Init()
		for _, testCase := range testCases {
			bucket := memoryStorage.calcBucket(testCase.input)
			require.Equal(t, testCase.expected, bucket)
		}
	})
}

func TestInitCounter(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		timeStamp := 1000000
		testCases := []struct {
			input    BucketInfo
			expected []BucketValue
		}{
			{
				input: BucketInfo{
					Interval: 2000,
					Limit:    4,
					Count:    3,
				},
				expected: []BucketValue{
					{
						Timestamp: timeStamp,
						Value:     1,
					},
					{
						Timestamp: timeStamp + 2000,
					},
					{
						Timestamp: timeStamp + 2000*2,
					},
				},
			},
		}

		memoryStorage := MemoryStorage{}
		memoryStorage.Init()
		for _, testCase := range testCases {
			counter := memoryStorage.initCounter(testCase.input, timeStamp)
			require.Equal(t, testCase.expected, counter)
		}
	})
}

func TestInc(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		item := struct {
			Login string
		}{
			Login: "Gogo",
		}

		limit := types.Limit{
			Interval: 2,
			Limit:    4,
		}

		memoryStorage := MemoryStorage{}
		memoryStorage.Init()
		require.True(t, memoryStorage.Inc(item, limit))
		require.True(t, memoryStorage.Inc(item, limit))
		time.Sleep(time.Millisecond * 1000)

		require.True(t, memoryStorage.Inc(item, limit))
		require.True(t, memoryStorage.Inc(item, limit))
		require.False(t, memoryStorage.Inc(item, limit))
		time.Sleep(time.Millisecond * 1500)

		require.True(t, memoryStorage.Inc(item, limit))
		require.True(t, memoryStorage.Inc(item, limit))
		require.False(t, memoryStorage.Inc(item, limit))
	})
}
