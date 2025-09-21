package util

import "time"

func NowUTCWithoutNanoSecond() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
