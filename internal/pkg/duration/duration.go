package duration

import (
	"time"

	"github.com/sirupsen/logrus"
)

func Duration(time1, time2 int64) string {
	t1 := time.Unix(time1, 0)
	t2 := time.Unix(time2, 0)
	diff := t2.Sub(t1)
	duration := time.Time{}.Add(diff).Format("15:04:05")
	return duration
}

func ParseTime(timeString string) (time.Time, error) {
	timeTime, err := time.Parse("2006-01-02, 15:04:05", timeString)
	if err != nil {
		logrus.Fatal(err)
	}
	return timeTime, nil
}
