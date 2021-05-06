package model

import "time"

func TimeAsEpochMicroseconds(t time.Time) uint64 {
	return uint64(t.UnixNano() / 1000)
}

func DurationAsMicroseconds(t time.Duration) uint64 {
	return uint64(t.Nanoseconds() / 1000)
}
