package difftimes

import (
	"fmt"
	"time"

	"bitbucket.org/pinang-och-go/pkg/v1/utils/constants"
)

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	now, _ := time.Parse("2006-01-02T15:04:05-0700", "2020-06-15T01:00:00+0700")
	fmt.Println("now", now)
	return now
}

func Diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func GetAge(formatTime, birthday, env string) (year int) {
	a, err := time.Parse(formatTime, birthday)
	if err != nil {
		return -1
	}

	b := time.Now()
	if env == "testing" {
		var mockClock realClock
		b = mockClock.Now()
	}

	year, _, _, _, _, _ = Diff(a, b)

	return
}

// ConvertToDestinedTimezone this function will return the time.Now() equal the timezone provided
// you can refer to TZ Database name here https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
func ConvertToDestinedTimezone(source time.Time, zone string) (time.Time, error) {
	jktTimezone, err := time.LoadLocation(zone)
	if err != nil {
		return source, err
	}
	return time.ParseInLocation(constants.ISOTimeLayout, source.Format(constants.ISOTimeLayout), jktTimezone)
}
