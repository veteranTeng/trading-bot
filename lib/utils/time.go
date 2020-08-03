package utils

import "time"

func GetLastWeekday(now time.Time) time.Time {
	prevDay := now.AddDate(0, 0, -1)
	for prevDay.Weekday() == time.Saturday || prevDay.Weekday() == time.Sunday {
		prevDay = prevDay.AddDate(0, 0, -1)
	}
	return prevDay
}

func GetMidnight(moment time.Time, loc *time.Location) time.Time {
	year, month, day := moment.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func GetStartOfWeek(moment time.Time, loc *time.Location) time.Time {
	current := moment
	for current.Weekday() != time.Monday {
		current = current.AddDate(0, 0, -1)
	}
	return GetMidnight(current, loc)
}

func GetStartOfMonth(moment time.Time, loc *time.Location) time.Time {
	year, month, _ := moment.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}

func GetMinuteBucket(moment time.Time, loc *time.Location, interval uint) time.Time {
	year, month, day := moment.Date()
	hour, minute, _ := moment.Clock()
	bucket := minute / int(interval) * int(interval)
	return time.Date(year, month, day, hour, bucket, 0, 0, loc)
}

func GetHourBucket(moment time.Time, loc *time.Location, interval uint) time.Time {
	year, month, day := moment.Date()
	hour, _, _ := moment.Clock()
	bucket := hour / int(interval) * int(interval)
	return time.Date(year, month, day, bucket, 0, 0, 0, loc)
}

func GetNYSELocation() *time.Location {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	return location
}

// Assuming location = "America/New_York"
func IsWithinNYSETradingHours(moment time.Time) bool {
	hour, minute, _ := moment.Clock()
	if hour == 9 {
		return minute >= 30
	}
	return hour > 9 && hour < 16
}
