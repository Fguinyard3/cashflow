package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// The timezone the database is set to - defaults to UTC
var DatabaseLocation, _ = time.LoadLocation("UTC")

// the amount of precision removed
var dateOnlyPrecision = 24 * time.Hour

// DateOnly represents a date-only value
type DateOnly struct {
	time.Time
}

// Scan satisfies the sql.scanner interface
func (d *DateOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		d.Time = truncateDateOnly(v)
	case string:
		if err := d.UnmarshalJSON([]byte(v)); err != nil {
			d = nil
			return fmt.Errorf("could not scan DateOnly value: %q", v)
		}
	default:
		return fmt.Errorf("could not scan DateOnly value of unknown type: %#v", value)
	}

	return nil
}

// Value satisfies the driver.Value interface
func (d DateOnly) Value() (driver.Value, error) {
	return truncateDateOnly(d.Time), nil // format just in case
}

// ensure the correct format is stored in the DB
func truncateDateOnly(t time.Time) time.Time {
	return t.In(DatabaseLocation).Truncate(dateOnlyPrecision)
}

const dateOnlyFormat = "2006-01-02"

var DateOnlyFormats = []string{dateOnlyFormat, "Jan 02, 2006", "2006-Jan-02"}

func (d DateOnly) String() string {
	return d.Format(dateOnlyFormat)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	// Add quotes
	return []byte(`"` + d.String() + `"`), nil
}

// DateOnly wrapper around the time.Date() function
func DateOnlyDate(year int, month time.Month, day int, loc *time.Location) *DateOnly {
	return &DateOnly{Time: truncateDateOnly(time.Date(year, month, day, 0, 0, 0, 0, loc))}
}

// Wrap an existing time.Time in a DateOnly
func NewDateOnly(t time.Time) *DateOnly {
	return &DateOnly{Time: truncateDateOnly(t)}
}

// Wraper around time.Now() in DateOnly format
func NowDateOnly() *DateOnly {
	return &DateOnly{Time: truncateDateOnly(NowSource())}
}

// UnmarshalJSON satisfies the json.Unmarshal interface
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var err error
	str := string(data)
	if str == "null" {
		d = nil
		return nil
	}
	// Only remove quotes if they are there
	if data[0] == '"' && data[len(str)-1] == '"' {
		// Remove quotes
		str = str[1 : len(str)-1]
	}
	// Try to parse with date only format
	var rt time.Time
	// Be lenient and try to parse from all the date only formats
	for _, layout := range DateOnlyFormats {
		if rt, err = time.Parse(layout, str); err == nil {
			rt = truncateDateOnly(rt)
			break
		}
	}
	if err != nil {
		// Be lenient and try to parse from all the other timestamp formats
		for _, layout := range TimestampFormats {
			if rt, err = time.Parse(layout, str); err == nil {
				rt = truncateDateOnly(rt)
				break
			}
		}
	}
	if err == nil {
		d.Time = rt
	} else {
		return fmt.Errorf("cannot parse date only value: %q", str)
	}
	return nil
}

// IsValid returns true if the pointer isn't nil and the value isnt zero
func (d *DateOnly) IsValid() bool {
	return d != nil && !d.IsZero()
}

func AddDenormalizedMonth(t time.Time, m int) time.Time {
	x := t.AddDate(0, m, 0)
	if d := x.Day(); d != t.Day() {
		return x.AddDate(0, 0, -d)
	}
	return x
}

// Returns a new DateOnly with the given number of months added
func NewDateOnlyFromAddedMonths(d *DateOnly, m int) *DateOnly {
	newTime := AddDenormalizedMonth(d.Time, m)
	return NewDateOnly(newTime)
}
