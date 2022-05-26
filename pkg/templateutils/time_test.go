package templateutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeNS_Constants(t *testing.T) {
	var ns timeNS

	assert.Equal(t, time.ANSIC, ns.FMT_ANSI())
	assert.Equal(t, time.RFC3339, ns.FMT_RFC3339())
	assert.Equal(t, time.RFC3339Nano, ns.FMT_RFC3339Nano())
	assert.Equal(t, time.RubyDate, ns.FMT_Ruby())
	assert.Equal(t, time.UnixDate, ns.FMT_Unix())
	assert.Equal(t, time.StampMilli, ns.FMT_Stamp())
}

func TestTimeNS_Duration(t *testing.T) {
	var ns timeNS
	for _, test := range []struct {
		gen func(i ...Number) (time.Duration, error)
		one time.Duration
	}{
		{ns.Nanosecond, time.Nanosecond},
		{ns.Microsecond, time.Microsecond},
		{ns.Millisecond, time.Millisecond},
		{ns.Second, time.Second},
		{ns.Minute, time.Minute},
		{ns.Hour, time.Hour},
		{ns.Day, 24 * time.Hour},
		{ns.Week, 7 * 24 * time.Hour},
	} {
		t.Run(test.one.String(), func(t *testing.T) {
			d, err := test.gen()
			assert.Equal(t, test.one, d)
			assert.NoError(t, err)

			d, err = test.gen(1)
			assert.Equal(t, test.one, d)
			assert.NoError(t, err)

			d, err = test.gen("1")
			assert.Equal(t, test.one, d)
			assert.NoError(t, err)

			d, err = test.gen(1.1)
			assert.Equal(t, test.one+test.one/10, d)
			assert.NoError(t, err)

			d, err = test.gen("1.1")
			assert.Equal(t, test.one+test.one/10, d)
			assert.NoError(t, err)
		})
	}

	t.Run("ParseDuration", func(t *testing.T) {
		d, err := ns.ParseDuration(1)
		assert.NoError(t, err)
		assert.Equal(t, time.Nanosecond, d)

		d, err = ns.ParseDuration("1")
		assert.NoError(t, err)
		assert.Equal(t, time.Nanosecond, d)

		d, err = ns.ParseDuration(1.1)
		assert.NoError(t, err)
		assert.Equal(t, time.Second+100*time.Millisecond, d)

		d, err = ns.ParseDuration("1.1")
		assert.NoError(t, err)
		assert.Equal(t, time.Second+100*time.Millisecond, d)

		d, err = ns.ParseDuration("1h")
		assert.NoError(t, err)
		assert.Equal(t, time.Hour, d)
	})

	t.Run("RoundDuration", func(t *testing.T) {
		d, err := ns.RoundDuration(time.Second, 501*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, time.Second, d)

		d, err = ns.RoundDuration(time.Second, 499*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), d)

		d, err = ns.FloorDuration(time.Second, 501*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), d)

		d, err = ns.FloorDuration(time.Second, 1501*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, time.Second, d)

		d, err = ns.CeilDuration(time.Second, time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, time.Second, d)

		d, err = ns.CeilDuration(time.Second, time.Second)
		assert.NoError(t, err)
		assert.Equal(t, time.Second, d)

		d, err = ns.CeilDuration(0, 0)
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), d)
	})
}

func TestTimeNS_Time(t *testing.T) {
	// Zone	NAME            STDOFF    RULES  FORMAT    [UNTIL]
	// Zone Pacific/Apia	12:33:04   -      LMT       1892 Jul  5
	//                      -11:26:56  -      LMT       1911
	//  	                -11:30	   -      -1130     1950
	//                      -11:00     WS	  -11/-10   2011 Dec 29 24:00
	//                      13:00      WS     +13/+14
	const TZ = "Pacific/Apia"

	tz, err := time.LoadLocation(TZ)
	if !assert.NoError(t, err) {
		return
	}

	var ns timeNS

	start := time.Date(1940, time.May, 10, 12, 31, 31, 50000, time.UTC)

	t.Run("Now", func(t *testing.T) {
		t.Parallel()

		ti, err := ns.Now()
		assert.NoError(t, err)
		assert.False(t, ti.IsZero())

		ti, err = ns.Now(TZ)
		assert.NoError(t, err)
		assert.False(t, ti.IsZero())
	})

	t.Run("ModTime", func(t *testing.T) {
		t.Parallel()

		ti, err := ns.Round(time.Hour, start)
		assert.NoError(t, err)
		assert.True(t,
			ti.Equal(time.Date(1940, time.May, 10, 13, 0, 0, 0, time.UTC)),
			"actual: %s\norigin: %s",
			ti.Format(time.RFC3339Nano),
			start.Format(time.RFC3339Nano),
		)

		ti, err = ns.Floor(time.Hour, start)
		assert.NoError(t, err)
		assert.True(t,
			ti.Equal(time.Date(1940, time.May, 10, 12, 0, 0, 0, time.UTC)),
			"actual: %s\norigin: %s",
			ti.Format(time.RFC3339Nano),
			start.Format(time.RFC3339Nano),
		)

		ti, err = ns.Ceil(time.Hour, start)
		assert.NoError(t, err)
		assert.True(t,
			ti.Equal(time.Date(1940, time.May, 10, 13, 0, 0, 0, time.UTC)),
			"actual: %s\norigin: %s",
			ti.Format(time.RFC3339Nano),
			start.Format(time.RFC3339Nano),
		)

		// start.In(tz) = 1940-05-10T01:01:31.00005-11:30

		ti, err = ns.Floor(time.Hour, start.In(tz))
		assert.NoError(t, err)
		assert.True(t,
			ti.Equal(time.Date(1940, time.May, 10, 1, 0, 0, 0, tz)),
			"actual: %s\norigin: %s",
			ti.Format(time.RFC3339Nano),
			start.In(tz).Format(time.RFC3339Nano),
		)

		ti, err = ns.Round(time.Hour, start.In(tz))
		assert.NoError(t, err)
		assert.True(t,
			ti.Equal(time.Date(1940, time.May, 10, 1, 0, 0, 0, tz)),
			"actual: %s\norigin: %s",
			ti.Format(time.RFC3339Nano),
			start.In(tz).Format(time.RFC3339Nano),
		)

		ti, err = ns.Ceil(time.Hour, start.In(tz))
		assert.NoError(t, err)
		assert.True(t,
			ti.Equal(time.Date(1940, time.May, 10, 2, 0, 0, 0, tz)),
			"actual: %s\norigin: %s",
			ti.Format(time.RFC3339Nano),
			start.In(tz).Format(time.RFC3339Nano),
		)
	})

	t.Run("Parse", func(t *testing.T) {
		t.Parallel()

		t.Run("time.Time", func(t *testing.T) {
			ti, err := ns.Parse(start)
			assert.NoError(t, err)
			assert.Equal(t, start, ti)

			ti, err = ns.Parse(time.ANSIC, start)
			assert.NoError(t, err)
			assert.Equal(t, start, ti)

			ti, err = ns.Parse(time.ANSIC, tz, start)
			assert.NoError(t, err)
			assert.Equal(t, start.In(tz), ti)

			ti, err = ns.Parse("", tz, start)
			assert.NoError(t, err)
			assert.Equal(t, start.In(tz), ti)
		})

		t.Run("String", func(t *testing.T) {
			ti, err := ns.Parse(start.Format(time.RFC3339))
			assert.NoError(t, err)
			expected, err := time.Parse(time.RFC3339, start.Format(time.RFC3339))
			assert.NoError(t, err)
			assert.Equal(t, expected, ti)

			ti, err = ns.Parse(time.RFC3339Nano, start.Format(time.RFC3339Nano))
			assert.NoError(t, err)
			expected, err = time.Parse(time.RFC3339Nano, start.Format(time.RFC3339Nano))
			assert.NoError(t, err)
			assert.Equal(t, expected, ti)

			ti, err = ns.Parse("", TZ, start.Format(time.RFC3339))
			assert.NoError(t, err)
			expected, err = time.Parse(time.RFC3339, start.Format(time.RFC3339))
			assert.NoError(t, err)
			assert.Equal(t, expected.In(tz), ti)
		})
	})

	t.Run("Format", func(t *testing.T) {
		t.Parallel()

		f, err := ns.Format(start)
		assert.NoError(t, err)
		assert.Equal(t, start.Format(time.RFC3339), f)

		f, err = ns.Format(time.ANSIC, start)
		assert.NoError(t, err)
		assert.Equal(t, start.Format(time.ANSIC), f)

		f, err = ns.Format(time.ANSIC, TZ, start)
		assert.NoError(t, err)
		assert.Equal(t, start.In(tz).Format(time.ANSIC), f)

		f, err = ns.Format("", TZ, start)
		assert.NoError(t, err)
		assert.Equal(t, start.In(tz).Format(time.RFC3339), f)
	})

	t.Run("Add", func(t *testing.T) {
		t.Parallel()

		ti, err := ns.Add(100, start)
		assert.NoError(t, err)
		assert.True(t, start.Add(100).Equal(ti))

		ti, err = ns.Add("100", start)
		assert.NoError(t, err)
		assert.True(t, start.Add(100).Equal(ti))

		ti, err = ns.Add(1.1, start)
		assert.NoError(t, err)
		assert.True(t, start.Add(time.Second+100*time.Millisecond).Equal(ti))

		ti, err = ns.Add("1.1", start)
		assert.NoError(t, err)
		assert.True(t, start.Add(time.Second+100*time.Millisecond).Equal(ti))

		ti, err = ns.Add(time.Second, start)
		assert.NoError(t, err)
		assert.True(t, start.Add(time.Second).Equal(ti))
	})

	t.Run("Since", func(t *testing.T) {
		t.Parallel()

		d, err := ns.Since(start)
		assert.NoError(t, err)
		assert.True(t, d > 0)

		d, err = ns.Since(start, start)
		assert.NoError(t, err)
		assert.True(t, d == 0)
	})

	t.Run("Until", func(t *testing.T) {
		t.Parallel()

		d, err := ns.Until(start)
		assert.NoError(t, err)
		assert.True(t, d < 0)

		d, err = ns.Until(start, start)
		assert.NoError(t, err)
		assert.True(t, d == 0)
	})
}
