package customtime

import (
	"fmt"
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
