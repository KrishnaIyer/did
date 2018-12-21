package day

import "time"

var now = time.Now()

func day(weekday time.Weekday) time.Time {
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for day.Weekday() != weekday {
		day = day.Add(-1 * time.Hour * 24)
	}
	return day
}

const (
	today     = "today"
	yesterday = "yesterday"
	monday    = "monday"
	tuesday   = "tuesday"
	wednesday = "wednesday"
	thursday  = "thursday"
	friday    = "friday"
	saturday  = "saturday"
	sunday    = "sunday"
)

var indicators = map[string]time.Time{
	today:     day(now.Weekday()),
	yesterday: day(now.Add(-1 * time.Hour * 24).Weekday()),
	monday:    day(time.Monday),
	tuesday:   day(time.Tuesday),
	wednesday: day(time.Wednesday),
	thursday:  day(time.Thursday),
	friday:    day(time.Friday),
	saturday:  day(time.Saturday),
	sunday:    day(time.Sunday),
}

// Indicators of days.
func Indicators() []string {
	return []string{
		today,
		yesterday,
		monday,
		tuesday,
		wednesday,
		thursday,
		friday,
		saturday,
		sunday,
	}
}

// GetMidnight gets the "midnight" moment for the given day indicator.
func GetMidnight(indicator string) (time.Time, bool) {
	t, ok := indicators[indicator]
	return t, ok
}
