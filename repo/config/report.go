package config

import (
	"fmt"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}

func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

type ReportInfo struct {
	Address        string
	DurationMin    Duration
	DurationMax    Duration
	RequestTimeout Duration
}

var defaultReportInfo = ReportInfo{
	Address:        "http://1.1.2.3/abc",
	DurationMin:    Duration{4 * time.Minute},
	DurationMax:    Duration{5 * time.Minute},
	RequestTimeout: Duration{5 * time.Second},
}
