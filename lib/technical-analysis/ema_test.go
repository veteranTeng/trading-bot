package analyze

import "testing"

func TestEMA(t *testing.T) {
	t.Run(
		"EMA not enough elements for any calculation",
		testEMAFunc(
			[]int64{10, 10, 10},
			5,
			make([]ValidMicro, 3),
		),
	)
	t.Run(
		"EMA enough elements for one calculation",
		testEMAFunc(
			[]int64{10, 10, 10, 10, 15},
			5,
			// Same as SMA for initial calc: (10 + 10 + 10 + 10 + 15) / 5 = 11
			[]ValidMicro{{}, {}, {}, {}, genValidMicro(11)},
		),
	)
	t.Run(
		"EMA enough elements for multiple calculations",
		testEMAFunc(
			[]int64{
				DollarsToMicros(14),
				DollarsToMicros(13),
				DollarsToMicros(14),
				DollarsToMicros(13),
				DollarsToMicros(12),
				DollarsToMicros(12),
				DollarsToMicros(11),
			},
			5,
			[]ValidMicro{
				{},
				{},
				{},
				{},
				// (14 + 13 + 14 + 13 + 12) / 5 = 13.2
				genValidMicro(13200000),
				// (12 - 13.2) * (2 / (5 + 1)) + 13.2 = 12.8
				genValidMicro(12800000),
				// (11 - 12.8) * (2 / (5 + 1)) + 12.8 = 12.2
				genValidMicro(12200000),
			},
		),
	)
}

func testEMAFunc(closingPrices []int64, interval int, expected []ValidMicro) func(*testing.T) {
	return func(t *testing.T) {
		actual := EMA(closingPrices, interval)
		if !eqValidCalcSlice(expected, actual) {
			t.Errorf("\nExpected: %v\nActual: %v", expected, actual)
		}
	}
}
