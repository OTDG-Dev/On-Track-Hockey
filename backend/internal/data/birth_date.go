package data

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

var ErrInvalidDateFormat = errors.New("invalid date format, expected: 2006-01-30")

type BirthDate time.Time

func (d *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return ErrInvalidDateFormat
	}
	*d = BirthDate(t)
	return nil
}

func (d BirthDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (d BirthDate) Year() int {
	return time.Time(d).Year()
}

func (d BirthDate) Value() (driver.Value, error) {
	return time.Time(d), nil
}
