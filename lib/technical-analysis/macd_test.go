package analyze

import "testing"

func TestStandardMACD(t *testing.T) {
	t.Run(
		"Standard MACD not enough elements for any calculation",
		testStandardMACDFunc(
			[]int64{10, 10, 10, 10},
			make([]ValidMicro, 4),
		),
	)
	t.Run(
		"Standard MACD enough elements for one calculation",
		testStandardMACDFunc(
			[]int64{
				DollarsToMicros(459.99),
				DollarsToMicros(448.85),
				DollarsToMicros(446.06),
				DollarsToMicros(450.81),
				DollarsToMicros(442.8),
				DollarsToMicros(448.97),
				DollarsToMicros(444.57),
				DollarsToMicros(441.4),
				DollarsToMicros(430.47),
				DollarsToMicros(420.05),
				DollarsToMicros(431.14),
				DollarsToMicros(425.66),
				DollarsToMicros(430.58),
				DollarsToMicros(431.72),
				DollarsToMicros(437.87),
				DollarsToMicros(428.43),
				DollarsToMicros(428.35),
				DollarsToMicros(432.5),
				DollarsToMicros(443.66),
				DollarsToMicros(455.72),
				DollarsToMicros(454.49),
				DollarsToMicros(452.08),
				DollarsToMicros(452.73),
				DollarsToMicros(461.91),
				DollarsToMicros(463.58),
				DollarsToMicros(461.14),
				DollarsToMicros(452.08),
				DollarsToMicros(442.66),
				DollarsToMicros(428.91),
				DollarsToMicros(429.79),
				DollarsToMicros(431.99),
				DollarsToMicros(427.72),
				DollarsToMicros(423.2),
				DollarsToMicros(426.21),
			},
			// Same as SMA for initial calc: (10 + 10 + 10 + 10 + 15) / 5 = 11
			[]ValidMicro{
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {},
				{-5108084, true},
			},
		),
	)
	t.Run(
		"Standard MACD enough elements for multiple calculations",
		testStandardMACDFunc(
			[]int64{
				DollarsToMicros(459.99),
				DollarsToMicros(448.85),
				DollarsToMicros(446.06),
				DollarsToMicros(450.81),
				DollarsToMicros(442.8),
				DollarsToMicros(448.97),
				DollarsToMicros(444.57),
				DollarsToMicros(441.4),
				DollarsToMicros(430.47),
				DollarsToMicros(420.05),
				DollarsToMicros(431.14),
				DollarsToMicros(425.66),
				DollarsToMicros(430.58),
				DollarsToMicros(431.72),
				DollarsToMicros(437.87),
				DollarsToMicros(428.43),
				DollarsToMicros(428.35),
				DollarsToMicros(432.5),
				DollarsToMicros(443.66),
				DollarsToMicros(455.72),
				DollarsToMicros(454.49),
				DollarsToMicros(452.08),
				DollarsToMicros(452.73),
				DollarsToMicros(461.91),
				DollarsToMicros(463.58),
				DollarsToMicros(461.14),
				DollarsToMicros(452.08),
				DollarsToMicros(442.66),
				DollarsToMicros(428.91),
				DollarsToMicros(429.79),
				DollarsToMicros(431.99),
				DollarsToMicros(427.72),
				DollarsToMicros(423.2),
				DollarsToMicros(426.21),
				DollarsToMicros(426.98),
				DollarsToMicros(435.69),
				DollarsToMicros(434.33),
				DollarsToMicros(429.8),
				DollarsToMicros(419.85),
				DollarsToMicros(426.24),
				DollarsToMicros(402.8),
				DollarsToMicros(392.05),
				DollarsToMicros(390.53),
				DollarsToMicros(398.67),
				DollarsToMicros(406.13),
				DollarsToMicros(405.46),
				DollarsToMicros(408.38),
				DollarsToMicros(417.2),
				DollarsToMicros(430.12),
				DollarsToMicros(442.78),
			},
			[]ValidMicro{
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {}, {}, {},
				{}, {}, {},
				{-5108084, true},
				{-4527496, true},
				{-3387777, true},
				{-2592274, true},
				{-2250615, true},
				{-2552088, true},
				{-2192264, true},
				{-3335498, true},
				{-4543441, true},
				{-5129228, true},
				{-4666182, true},
				{-3602783, true},
				{-2729465, true},
				{-1785740, true},
				{-466764, true},
				{1280988, true},
				{3186354, true},
			},
		),
	)
}

func testStandardMACDFunc(closingPrices []int64, expected []ValidMicro) func(*testing.T) {
	return func(t *testing.T) {
		actual := StandardMACD(closingPrices)
		if !eqValidCalcSlice(expected, actual) {
			t.Errorf("\nExpected: %v\nActual: %v", expected, actual)
		}
	}
}
