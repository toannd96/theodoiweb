package duration

import (
	"reflect"
	"testing"
	"time"
)

const testTimeLayout = "2006-01-02, 15:04:05"

func TestDuration(t *testing.T) {
	type args struct {
		time1 int64
		time2 int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return duration between two timestamp",
			args: args{
				time1: 1657091090,
				time2: 1657091095,
			},
			want: "00:00:05",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.time1, tt.args.time2); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	timeWant, _ := time.Parse(testTimeLayout, "2022-07-06, 00:25:23")
	type args struct {
		timeString string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "should convert time string to time time",
			args: args{
				timeString: "2022-07-06, 00:25:23",
			},
			want:    timeWant,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.args.timeString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
