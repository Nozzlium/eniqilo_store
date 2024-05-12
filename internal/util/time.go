package util

import "time"

func Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}

func ToISO8601(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.999Z")
}
