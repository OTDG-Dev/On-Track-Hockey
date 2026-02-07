package data

import (
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

type TimeOnIce struct {
	time.Duration
}

// need to write Decode still
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
