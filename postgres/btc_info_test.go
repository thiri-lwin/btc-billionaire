package postgres

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// eg. if freq=30 and hour=2, btc is sent every 30 mins for 2 hours at the given timezone
// if both freq and hour is zero, assume that only one row will be inserted for testing
func insertBTCInfo(t *testing.T, freq int, hour int) {
	noOfTest := 1
	if freq > 0 && hour > 0 {
		noOfTest = (hour * 60) / freq
	}

	created := "2020-11-05T13:00:00+03:00"
	for noOfTest > 0 {
		createdAt, _ := time.Parse(time.RFC3339, created)
		createdAt.Add(time.Minute * time.Duration(freq))
		_, offset := createdAt.Zone()
		btcInfo := BTCInfo{
			Amount:  decimal.NewFromInt(10),
			Created: createdAt,
			Offset:  offset,
		}
		result, err := dbStorage.InsertBTCInfo(btcInfo)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, createdAt.UTC(), result.Created)
		require.NotZero(t, result.ID)
		noOfTest--
	}

}

func TestInsertBTCInfo(t *testing.T) {
	insertBTCInfo(t, 0, 0)
}

func TestGetBTCHistory(t *testing.T) {
	start := "2020-11-05T12:20:05+03:00"
	end := "2022-12-05T13:00:00+03:00"
	startDate, _ := time.Parse(time.RFC3339, start)
	endDate, _ := time.Parse(time.RFC3339, end)
	t.Run("Send BTC every 30 mins for 2 hours", func(t *testing.T) {
		insertBTCInfo(t, 30, 2)
		result, err := dbStorage.GetBTCHistory(startDate, endDate)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotZero(t, len(result))
		require.Equal(t, 2, len(result))
	})
}
