// Package bigdur parses large durations.
// In addition to units supported by the standard library,
// It supports:
//
// - d for 24 hours
//
// - mo for 30 days
//
// - y for 12 months
package bigdur

import (
	"bytes"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Duration represents an interval of time.
type Duration uint64

// Parse parses a Duration.
func Parse(str string) (Duration, error) {
	var (
		dur Duration

		parsingNum = true
		numStart   int
		numEnd     int
		i          int
		r          rune
	)

	finishPart := func() error {
		c, err := strconv.ParseFloat(str[numStart:numEnd], 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse coefficient")
		}
		if c < 0 {
			return errors.New("negative durations are invalid")
		}

		u, err := parseDurationUnit(str[numEnd:i])
		if err != nil {
			return errors.Wrapf(err, "failed to parse unit %q", str[numEnd:i])
		}

		dur += Duration(c * float64(u.Multiplier))

		numStart, numEnd = i, i
		return nil
	}

	for i, r = range str {
		// fmt.Printf("parsing %c, i %v\n", r, i)
		// check if numeric
		if (r >= '0' && r <= '9') || r == '.' {
			if !parsingNum {
				// fmt.Printf("1\n")
				if err := finishPart(); err != nil {
					return 0, err
				}
			}
			parsingNum = true
			continue
		}
		if parsingNum {
			numEnd = i
			parsingNum = false
		}
	}

	if numEnd > numStart {
		i = len(str)
		if err := finishPart(); err != nil {
			return 0, err
		}
	}

	return dur, nil
}

// Duration converts into a standard library time.Duration.
func (d Duration) Std() time.Duration {
	return time.Duration(time.Second * time.Duration(d))
}

// String returns a human-readable representation of Duration.
func (d Duration) String() string {
	out := &bytes.Buffer{}

	var (
		co Duration
	)

	if d == Duration(0) {
		return "0s"
	}

	for i := len(durationUnits) - 1; i >= 0 && d > 0; i-- {
		// fmt.Printf("d:%v %v:%v \n", uint64(d), durationUnits[i].Unit, uint64(durationUnits[i].Multiplier))
		if d >= durationUnits[i].Multiplier {
			co = d / durationUnits[i].Multiplier
			out.WriteString(strconv.FormatUint(uint64(co), 10))
			out.WriteString(durationUnits[i].Unit)
			d -= co * durationUnits[i].Multiplier
		}
	}

	return out.String()
}

type durationUnit struct {
	Unit       string
	Multiplier Duration
}

// Duration units
var (
	durationUnits = [...]durationUnit{
		{Unit: "s", Multiplier: 1},
		{Unit: "m", Multiplier: 60},
		{Unit: "h", Multiplier: 60 * 60},
		{Unit: "d", Multiplier: 60 * 60 * 24},
		{Unit: "w", Multiplier: 60 * 60 * 24 * 7},
		{Unit: "mo", Multiplier: 60 * 60 * 24 * 7 * 30},
		{Unit: "y", Multiplier: 60 * 60 * 24 * 7 * 30 * 12},
	}
)

func parseDurationUnit(unit string) (durationUnit, error) {
	for _, du := range durationUnits {
		if du.Unit == unit {
			return du, nil
		}
	}
	return durationUnit{}, errors.New("no unit found")
}
