package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidDateFormat = errors.New("invalid date format, expected: 2006-01-30")

type BirthDate struct {
	time.Time
}

func (d *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return ErrInvalidDateFormat
	}
	d.Time = t
	return nil
}

func (d BirthDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d BirthDate) Year() int {
	return d.Time.Year()
}

var ErrInvalidTimeOnIceFormat = errors.New("invalid time ice format, expected: MM:SS")

type TimeOnIce struct {
	time.Duration
}

func (t TimeOnIce) MarshalJSON() ([]byte, error) {
	if t.Duration <= 0 {
		return []byte(`"00:00"`), nil
	}

	totalSeconds := int(t.Seconds())
	mins := totalSeconds / 60
	secs := totalSeconds % 60

	var b []byte
	return fmt.Appendf(b, `"%02d:%02d"`, mins, secs), nil
}

// WIP
func (t *TimeOnIce) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidTimeOnIceFormat
	}

	parts := strings.Split(unquotedJSONValue, ":")

	_ = parts
	return nil
}
