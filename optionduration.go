package retryer

import (
	"fmt"
	"time"
)

type OptionDuration time.Duration

func (od *OptionDuration) UnmarshalText(b []byte) (err error) {
	var d time.Duration
	d, err = time.ParseDuration(string(b))
	*od = OptionDuration(d)
	return
}

func (od *OptionDuration) MarshalText() ([]byte, error) {
	return []byte(od.String()), nil
}

func (od *OptionDuration) String() string {
	d := time.Duration(*od)
	return fmt.Sprintf("%v", d)
}
