package timeutil

import "time"

func ParseDuration(s string) (time.Duration, error) {

	return time.ParseDuration(s)
}
