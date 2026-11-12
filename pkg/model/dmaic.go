package model

import (
	"regexp"
	"strconv"
)

// DMAIC readiness thresholds, shared by lint and the DMAIC export so the two
// agree on what "enough data" means.
const (
	// DMAICMinPoints is the fewest measurements that let a DMAIC effort set
	// control limits: statistical process control asks for 20 to 25 points
	// before the limits mean anything. Lint asks for this many a year; the
	// export computes control limits only from this many observations.
	DMAICMinPoints = 20
	// DMAICMinResponsesPerPoint is the fewest survey responses behind one
	// measurement before its value stops swinging with chance. An NPS needs
	// about 100 for a margin near ±10 points; 30 is the floor.
	DMAICMinResponsesPerPoint = 30
)

var isoDuration = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

// TimesPerYear is how many times a year an ISO 8601 duration repeats, using the
// average year and month. It reports false for an empty or unreadable
// duration and for one of zero length.
func TimesPerYear(iso string) (float64, bool) {
	m := isoDuration.FindStringSubmatch(iso)
	if m == nil {
		return 0, false
	}
	const daysPerYear = 365.2425
	unit := []float64{daysPerYear, daysPerYear / 12, 7, 1, 1.0 / 24, 1.0 / 1440, 1.0 / 86400}
	days := 0.0
	for i, part := range m[1:] {
		if part == "" {
			continue
		}
		v, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return 0, false // the pattern admits only numbers; unreachable
		}
		days += v * unit[i]
	}
	if days == 0 {
		return 0, false
	}
	return daysPerYear / days, true
}
