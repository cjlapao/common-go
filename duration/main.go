// Package duration provides a partial implementation of ISO8601 durations. (no months)
package duration

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

var (
	// ErrBadFormat is returned when parsing fails
	ErrBadFormat = errors.New("bad format string")

	// ErrNoMonth is raised when a month is in the format string
	ErrNoMonth = errors.New("no months allowed")

	tmpl = template.Must(template.New("duration").Parse(`P{{if .Years}}{{.Years}}Y{{end}}{{if .Weeks}}{{.Weeks}}W{{end}}{{if .Days}}{{.Days}}D{{end}}{{if .HasTimePart}}T{{end }}{{if .Hours}}{{.Hours}}H{{end}}{{if .Minutes}}{{.Minutes}}M{{end}}{{if .Seconds}}{{.Seconds}}S{{end}}`))

	full = regexp.MustCompile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+(?:\.\d+))S)?)?`)
	week = regexp.MustCompile(`P((?P<week>\d+)W)`)
)

type Duration struct {
	Years        int
	Weeks        int
	Days         int
	Hours        int
	Minutes      int
	Seconds      int
	MilliSeconds int
}

func FromString(dur string) (*Duration, error) {
	var (
		match []string
		re    *regexp.Regexp
	)

	if week.MatchString(dur) {
		match = week.FindStringSubmatch(dur)
		re = week
	} else if full.MatchString(dur) {
		match = full.FindStringSubmatch(dur)
		re = full
	} else {
		return nil, ErrBadFormat
	}

	d := &Duration{}

	for i, name := range re.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		val, err := strconv.ParseFloat(part, 10)
		if err != nil {
			return nil, err
		}
		switch name {
		case "year":
			d.Years = int(val)
		case "month":
			return nil, ErrNoMonth
		case "week":
			d.Weeks = int(val)
		case "day":
			d.Days = int(val)
		case "hour":
			c := time.Duration(val) * time.Hour
			d.Hours = int(c.Hours())
		case "minute":
			c := time.Duration(val) * time.Minute
			d.Minutes = int(c.Minutes())
		case "second":
			s, milli := math.Modf(val)
			d.Seconds = int(s)
			d.MilliSeconds = int(milli * 1000)
		default:
			return nil, fmt.Errorf("unknown field %s", name)
		}
	}

	return d, nil
}

// String prints out the value passed in. It's not strictly according to the
// ISO spec, but it's pretty close. It would need to disallow weeks mingling with
// other units.
func (d *Duration) String() string {
	var s bytes.Buffer

	d.normalize()

	err := tmpl.Execute(&s, d)
	if err != nil {
		panic(err)
	}

	return s.String()
}

func (d *Duration) normalize() {
	msToS := 1000
	StoM := 60
	MtoH := 60
	HtoD := 24
	DtoW := 7
	if d.MilliSeconds >= msToS {
		d.Seconds += d.MilliSeconds / msToS
		d.MilliSeconds %= msToS
	}
	if d.Seconds >= StoM {
		d.Minutes += d.Seconds / StoM
		d.Seconds %= StoM
	}
	if d.Minutes >= MtoH {
		d.Hours += d.Minutes / MtoH
		d.Minutes %= MtoH
	}
	if d.Hours >= HtoD {
		d.Days += d.Hours / HtoD
		d.Hours %= HtoD
	}
	if d.Days >= DtoW {
		d.Weeks += d.Days / DtoW
		d.Days %= DtoW
	}
	// a month is not always 30 days, so we don't normalize that
	// a month is not always 4 weeks, so we don't normalize that
	// a year is not always 52 weeks, so we don't normalize that
}

func (d *Duration) HasTimePart() bool {
	return d.Hours != 0 || d.Minutes != 0 || d.Seconds != 0
}

func (d *Duration) ToDuration() time.Duration {
	day := time.Hour * 24
	year := day * 365

	tot := time.Duration(0)

	tot += year * time.Duration(d.Years)
	tot += day * 7 * time.Duration(d.Weeks)
	tot += day * time.Duration(d.Days)
	tot += time.Hour * time.Duration(d.Hours)
	tot += time.Minute * time.Duration(d.Minutes)
	tot += time.Second * time.Duration(d.Seconds)

	return tot
}
