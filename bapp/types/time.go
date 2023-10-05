package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var NowSource = time.Now

var TimestampPrecision = time.Microsecond

type Timestamp struct {
	time.Time `gorm:"type:timestamptz"`
}

var TimestampFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC850,
	time.RFC822Z,
	time.RFC822,
	time.RubyDate,
	time.UnixDate,
	time.ANSIC,
	"2006-01-02 15:04:05-07:00", // sqlite datetime string format
}

// satisfy the sql.scanner interface
func (t *Timestamp) Scan(value interface{}) error {
	var rt time.Time

	switch v := value.(type) {
	case time.Time:
		rt = v
	case string:
		// Try to parse it in all the layouts
		var err error
		for _, layout := range TimestampFormats {
			if rt, err = time.Parse(layout, v); err == nil {
				break
			}
		}
		if err != nil {
			return fmt.Errorf("cannot convert timestamp value: %v", value)
		}
	default:
		return fmt.Errorf("cannot convert timestamp value: %#v", value)
	}

	*t = Timestamp{Time: truncateTimestamp(rt)}
	return nil
}

// satisfies the driver.Value interface
func (t Timestamp) Value() (driver.Value, error) {
	return truncateTimestamp(t.Time), nil
}

// Now wrapper around the time.Now() function
func NowTimestamp() *Timestamp {
	return &Timestamp{Time: truncateTimestamp(NowSource())}
}

// Date wrapper around the time.Date() function
func TimestampDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *Timestamp {
	return &Timestamp{Time: truncateTimestamp(time.Date(year, month, day, hour, min, sec, nsec, loc))}
}

// Wrap an existing time.Time in a Timestamp
func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{Time: truncateTimestamp(t)}
}

func truncateTimestamp(t time.Time) time.Time {
	return t.In(DatabaseLocation).Truncate(TimestampPrecision)
}

func (t *Timestamp) IsValid() bool {
	return t != nil && !t.IsZero()
}
