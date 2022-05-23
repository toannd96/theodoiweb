package duration

import "time"

func Duration(time1, time2 int64) string {
	t1 := time.Unix(time1, 0)
	t2 := time.Unix(time2, 0)
	diff := t2.Sub(t1)
	duration := time.Time{}.Add(diff).Format("15:04:05")
	return duration
}
