package templateutils

import (
	"math"
	"time"
)

// timeNS for time
type timeNS struct{}

/*

	Constants

*/

func (timeNS) FMT_RFC3339() string     { return time.RFC3339 }
func (timeNS) FMT_RFC3339Nano() string { return time.RFC3339Nano }
func (timeNS) FMT_Unix() string        { return time.UnixDate }
func (timeNS) FMT_Ruby() string        { return time.RubyDate }
func (timeNS) FMT_ANSI() string        { return time.ANSIC }
func (timeNS) FMT_Stamp() string       { return time.StampMilli }
func (timeNS) FMT_Date() string        { return "2006-01-02" }
func (timeNS) FMT_Clock() string       { return "15:04:05" }
func (timeNS) FMT_DateTime() string    { return "2006-01-02T15:04:05" }

/*

	Duration related

*/

// ParseDuration parses a duration like 1h2m, 1, -1, 1.2
func (timeNS) ParseDuration(n String) (time.Duration, error) { return parseDuration(n) }

// FloorDuration floor duration x to mutiple of dur
func (timeNS) FloorDuration(dur, x Duration) (time.Duration, error) {
	return modDuration(dur, x, mod_Floor)
}

// RoundDuration round duration x to mutiple of dur
func (timeNS) RoundDuration(dur, x Duration) (time.Duration, error) {
	return modDuration(dur, x, mod_Round)
}

// CeilDuration floor duration x to mutiple of dur
func (timeNS) CeilDuration(dur, x Duration) (time.Duration, error) {
	return modDuration(dur, x, mod_Ceil)
}

func modDuration(a, x Duration, action modAction) (ret time.Duration, err error) {
	da, err := parseDuration(a)
	if err != nil {
		return
	}

	dx, err := parseDuration(x)
	if err != nil {
		return
	}

	return ModDuration(da, dx, action), nil
}

// Nanosecond get a duration of n ns
func (timeNS) Nanosecond(n ...Number) (time.Duration, error) { return mulDuration(time.Nanosecond, n) }

// Microsecond get a duration of n microseonds
func (timeNS) Microsecond(n ...Number) (time.Duration, error) {
	return mulDuration(time.Microsecond, n)
}

// Millisecond gets a duration of n ms
func (timeNS) Millisecond(n ...Number) (time.Duration, error) {
	return mulDuration(time.Millisecond, n)
}

// Second gets a duration of n s
func (timeNS) Second(n ...Number) (time.Duration, error) { return mulDuration(time.Second, n) }

// Minute gets a duration of n m
func (timeNS) Minute(n ...Number) (time.Duration, error) { return mulDuration(time.Minute, n) }

// Hour gets a duration of n h
func (timeNS) Hour(n ...Number) (time.Duration, error) { return mulDuration(time.Hour, n) }

// Day gets a duration of n days
func (timeNS) Day(n ...Number) (time.Duration, error) { return mulDuration(time.Hour*24, n) }

// Week gets a duration of n weeks
func (timeNS) Week(n ...Number) (time.Duration, error) { return mulDuration(time.Hour*24*7, n) }

func mulDuration(x time.Duration, dur []Number) (_ time.Duration, err error) {
	var (
		d       uint64
		isFloat bool
	)
	if len(dur) == 0 {
		d = 1
	} else {
		d, isFloat, err = parseNumber(dur[0])
		if err != nil {
			return
		}
	}

	if isFloat {
		return time.Duration(math.Float64frombits(d) * float64(x)), nil
	}

	return x * time.Duration(d), nil
}

/*

	Timezone related

*/

// ZoneName gets timezone name of time
//
// ZoneName(): get local timezone name
//
// ZoneName(t Time): get timezone name of t
//
// ZoneName(layout String, t Time): parse t in layout and get timezone name of it
func (timeNS) ZoneName(args ...any) (name string, err error) {
	switch n := len(args); n {
	case 0:
		// NOTE: do not use zero time as timezone name changes
		name, _ = time.Now().In(time.Local).Zone()
		return
	case 1:
		var t time.Time
		t, err = toTimeDefault(args[0])
		if err != nil {
			return
		}

		name, _ = t.Zone()
		return
	default:
		var (
			layout string
			t      time.Time
		)

		layout, err = toString(args[0])
		if err != nil {
			return
		}

		t, err = parseTime(layout, args[n-1], time.Local)
		if err != nil {
			return
		}

		name, _ = t.Zone()
		return
	}
}

// ZoneOffset get offset to UTC timezone
//
// ZoneOffset(): get local timezone offset
//
// ZoneOffset(t Time): get timezone offset of t
//
// ZoneOffset(layout String, t Time): parse t in layout and get timezone offset of it
func (timeNS) ZoneOffset(args ...any) (offset int, err error) {
	switch n := len(args); n {
	case 0:
		// NOTE: do not use zero time as timezone offset changes
		_, offset = time.Now().In(time.Local).Zone()
		return
	case 1:
		var t time.Time
		t, err = toTimeDefault(args[0])
		if err != nil {
			return
		}

		_, offset = t.Zone()
		return
	default:
		var (
			layout string
			t      time.Time
		)

		layout, err = toString(args[0])
		if err != nil {
			return
		}

		t, err = parseTime(layout, args[n-1], time.Local)
		if err != nil {
			return
		}

		_, offset = t.Zone()
		return
	}
}

/*

	Time related

*/

// Now get the current time
//
// Now(): get local current time
//
// Now(tz String): get current time in timezone tz
func (timeNS) Now(args ...String) (t time.Time, err error) {
	switch len(args) {
	case 0:
		return time.Now(), nil
	default:
		var (
			tz       string
			location *time.Location
		)

		tz, err = toString(args[0])
		if err != nil {
			return
		}

		location, err = time.LoadLocation(tz)
		if err != nil {
			return
		}

		return time.Now().In(location), nil
	}
}

// Round round time to multiple of duration dur
// unlike golang time.Duration.{Round, Truncate}, it zeros minutes in timezone of the time t instead of
// of in UTC timezone
func (timeNS) Round(dur Duration, t Time) (time.Time, error) {
	return modTime(dur, t, mod_Round)
}

// Floor is like Round but does floor unconditionally
func (timeNS) Floor(dur Duration, t Time) (time.Time, error) {
	return modTime(dur, t, mod_Floor)
}

// Ceil is like Round but does ceil unconditionally
func (timeNS) Ceil(dur Duration, t Time) (time.Time, error) {
	return modTime(dur, t, mod_Ceil)
}

func modTime(dur Duration, t Time, action modAction) (ret time.Time, err error) {
	d, err := parseDuration(dur)
	if err != nil {
		return
	}

	ti, err := toTimeDefault(t)
	if err != nil {
		return
	}

	return ModTime(d, ti, action), nil
}

// Parse last argument as time.Time, which can be a
// - ~string/~[]~byte (formatted time value)
// - number (seconds since unix epoch)
//
// Parse(t Time): parse Time t
// 		for string or bytes, in RFC3339Nano format (which accepts RFC3339)
//
// Parse(layout String, t Time): parse Time t
// 		for string or bytes, in specified layout
//
// Parse(layout, zone String, t Time): parse Time t with specified format
// 		when layout is not empty/nil, and t is string or bytes, then parse in specified layout
func (timeNS) Parse(args ...any) (ret time.Time, err error) {
	var (
		layout string
		loc    string
		t      Time
	)

	switch n := len(args); n {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		t = args[0]
	case 2:
		layout, err = toString(args[0])
		if err != nil {
			return
		}

		t = args[1]
	default:
		layout, err = toString(args[0])
		if err != nil {
			return
		}

		loc, err = toString(args[1])
		if err != nil {
			return
		}

		t = args[n-1]
	}

	var location *time.Location
	if len(loc) != 0 {
		location, err = time.LoadLocation(loc)
		if err != nil {
			return
		}
	}

	if len(layout) == 0 {
		layout = time.RFC3339Nano
	}

	return parseTime(layout, t, location)
}

// Format last argument (usually a time.Time) in time layout
//
// Format(t Time): format Time t in RFC3339 format (NOTE: not RFC3339Nano)
//
// Format(layout String, t Time): format Time t in specified layout
//
// Format(layout, loc String, t Time): format Time t in location loc wtih format, when format is empty, fallback to RFC3339
func (timeNS) Format(args ...any) (_ string, err error) {
	var (
		layout string
		loc    string
		t      Time
	)

	switch len(args) {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		t = args[0]
	case 2:
		layout, err = toString(args[0])
		if err != nil {
			return
		}

		t = args[1]
	default:
		layout, err = toString(args[0])
		if err != nil {
			return
		}

		loc, err = toString(args[1])
		if err != nil {
			return
		}

		t = args[2]
	}

	ti, err := toTimeDefault(t)
	if err != nil {
		return
	}

	if len(layout) == 0 {
		layout = time.RFC3339
	}

	if len(loc) != 0 {
		var location *time.Location

		location, err = time.LoadLocation(loc)
		if err != nil {
			return
		}

		return ti.In(location).Format(layout), nil
	}

	return ti.Format(layout), nil
}

// Add duration to time t
func (timeNS) Add(dur Duration, t Time) (ret time.Time, err error) {
	d, err := parseDuration(dur)
	if err != nil {
		return
	}

	ret, err = toTimeDefault(t)
	if err != nil {
		return
	}

	return ret.Add(d), nil
}

// Since get the duration since start time
//
// Since(start Time): return the time elapsed since start (until time.Now)
//
// Since(now, start Time): return time elapsed from start to end
func (timeNS) Since(args ...Time) (time.Duration, error) {
	switch len(args) {
	case 0:
		return 0, errAtLeastOneArgGotZero
	case 1: // Since(start)
		start, err := toTimeDefault(args[0])
		if err != nil {
			return 0, err
		}

		return time.Since(start), nil

	default: // Since(now, start)
		now, err := toTimeDefault(args[0])
		if err != nil {
			return 0, err
		}

		start, err := toTimeDefault(args[1])
		if err != nil {
			return 0, err
		}

		return now.Sub(start), nil
	}
}

// Until get the duration until end time
//
// Until(end Time): return duration from time.Now to end
// Until(now, end Time): return duration from now to end
func (timeNS) Until(args ...Time) (time.Duration, error) {
	switch len(args) {
	case 0:
		return 0, errAtLeastOneArgGotZero
	case 1:
		end, err := toTimeDefault(args[0])
		if err != nil {
			return 0, err
		}

		return time.Until(end), nil
	default:
		now, err := toTimeDefault(args[0])
		if err != nil {
			return 0, err
		}

		end, err := toTimeDefault(args[1])
		if err != nil {
			return 0, err
		}

		return end.Sub(now), nil
	}
}
