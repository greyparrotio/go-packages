package customtime

import (
	"fmt"
	"strings"
	"time"
)

type CustomTime time.Time

func (customTime CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(customTime)

	if t.IsZero() {
		return []byte(`""`), nil
	} else {
		return fmt.Appendf(nil, `"%s"`, t.Format(time.RFC3339)), nil
	}
}

func (customTime CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}

	s = strings.Trim(s, `"`)
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	customTime.Time = CustomTime(parsedTime)
	return nil
}
