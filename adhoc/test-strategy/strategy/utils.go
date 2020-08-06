package strategy

import (
	"tradingbot/lib/candle"
	analyze "tradingbot/lib/technical-analysis"
)

func groupCandlesByDate(candles []candle.Candle) ([]string, map[string][]candle.Candle) {
	dates := make([]string, 0)
	candlesByDate := make(map[string][]candle.Candle)
	for _, c := range candles {
		date := c.OpenedAt.Format("2006-01-02")
		groupedCandles, ok := candlesByDate[date]
		if !ok {
			dates = append(dates, date)
		}
		candlesByDate[date] = append(groupedCandles, c)
	}

	return dates, candlesByDate
}

func genInitialPortfolio(capitalMicros, priceMicros int64) Portfolio {
	shares := capitalMicros / priceMicros
	return Portfolio{
		Date:               "initial",
		Cash:               capitalMicros - priceMicros*shares,
		SharesHeld:         shares,
		EndOfDayValue:      capitalMicros,
		DailyChange:        0,
		DailyPercentChange: 0,
		AllTimePerformance: 0,
	}
}

func calcPercentChange(prevMicros int64, currentMicros int64) float64 {
	ratio := candle.MicrosToDollars(currentMicros) / candle.MicrosToDollars(prevMicros)
	return analyze.RoundToTwoDecimals(ratio*100 - 100)
}
