package counter

import (
	"fmt"
	"testing"
	"time"
)

func simpleTime(s string) time.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func sparseBuckets(buckets [3600]uint64) string {
	result := ""
	for i := range buckets {
		if buckets[i] != 0 {
			minutes := i / 60
			seconds := i % 60
			result += fmt.Sprintf("Bucket %4d (%02d:%02d): %d\n", i, minutes, seconds, buckets[i])
		}
	}
	return result
}

func TestCounter(t *testing.T) {
	testCases := []struct {
		name       string
		skip       bool
		countTimes []time.Time
		endTime    time.Time
		wantSecond uint64
		wantMinute uint64
		wantHour   uint64
	}{
		{
			name:       "zero",
			countTimes: []time.Time{},
			endTime:    simpleTime("12:01:00"),
			wantSecond: 0,
			wantMinute: 0,
			wantHour:   0,
		},
		{
			name: "one",
			countTimes: []time.Time{
				simpleTime("12:01:00"),
			},
			endTime:    simpleTime("12:01:00").Add(100 * time.Millisecond),
			wantSecond: 1,
			wantMinute: 1,
			wantHour:   1,
		},
		{
			name: "two-for-three",
			countTimes: []time.Time{
				simpleTime("12:01:00"),
				simpleTime("12:01:00"),
				simpleTime("12:01:01"),
				simpleTime("12:01:01"),
				simpleTime("12:01:02"),
				simpleTime("12:01:02"),
			},
			endTime:    simpleTime("12:01:02").Add(100 * time.Millisecond),
			wantSecond: 2,
			wantMinute: 6,
			wantHour:   6,
		},
		{
			name: "ran_out_of_minute",
			countTimes: []time.Time{
				simpleTime("12:01:00"),
				simpleTime("12:01:59"),
				simpleTime("12:02:00"),
			},
			endTime:    simpleTime("12:02:00").Add(100 * time.Millisecond),
			wantSecond: 1,
			wantMinute: 2,
			wantHour:   3,
		},
		{
			name: "once_in_30sec",
			countTimes: []time.Time{
				simpleTime("12:01:00"),
			},
			endTime:    simpleTime("12:01:30"),
			wantSecond: 0,
			wantMinute: 1,
			wantHour:   1,
		},
		{
			name: "once_in_2min",
			countTimes: []time.Time{
				simpleTime("12:01:00"),
			},
			endTime:    simpleTime("12:03:00"),
			wantSecond: 0,
			wantMinute: 0,
			wantHour:   1,
		},
		{
			name: "once_in_2hour",
			countTimes: []time.Time{
				simpleTime("12:01:00"),
			},
			endTime:    simpleTime("14:01:00"),
			wantSecond: 0,
			wantMinute: 0,
			wantHour:   0,
		},
		{
			name: "30_min_ticker_nojitter",
			countTimes: []time.Time{
				simpleTime("12:00:00"),
				simpleTime("12:00:00"),
				simpleTime("12:00:00"),
				simpleTime("12:00:00"),
				simpleTime("12:30:00"),
				simpleTime("12:30:00"),
				simpleTime("13:00:00"),
				simpleTime("13:30:00"),
			},
			endTime:    simpleTime("13:30:02"),
			wantSecond: 0,
			wantMinute: 1,
			wantHour:   2,
		},
		{
			name: "30_min_ticker_jitter",
			countTimes: []time.Time{
				simpleTime("12:00:01"),
				simpleTime("12:30:02"),
				simpleTime("13:00:03"),
				simpleTime("13:30:04"),
			},
			endTime:    simpleTime("13:30:06"),
			wantSecond: 0,
			wantMinute: 1,
			wantHour:   2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}
			c := New()
			for _, t := range tc.countTimes {
				c.Count(t)
			}

			if gotSecond := c.GetSecond(tc.endTime); gotSecond != tc.wantSecond {
				t.Errorf("c.GetSecond() is %v; want %v", gotSecond, tc.wantSecond)
			}
			if gotMinute := c.GetMinute(tc.endTime); gotMinute != tc.wantMinute {
				t.Errorf("c.GetMinute() is %v; want %v", gotMinute, tc.wantMinute)
			}
			if gotHour := c.GetHour(tc.endTime); gotHour != tc.wantHour {
				t.Errorf("c.GetHour() is %v; want %v", gotHour, tc.wantHour)
			}

			t.Log("Buckets:", sparseBuckets(c.buckets))
		})
	}
}
