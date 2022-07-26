package security

import (
	"testing"
)

func TestDoPasswordsMatch(t *testing.T) {
	type args struct {
		hashedPassword string
		currPassword   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test do password match",
			args: args{
				hashedPassword: "$2a$04$7m3nNCHR1JrI19jy/ZeLY.5F3ZVXd2Cac.EVj0kEeoQ2WxSVQYOhu",
				currPassword:   "12345678",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DoPasswordsMatch(tt.args.hashedPassword, tt.args.currPassword); got != tt.want {
				t.Errorf("DoPasswordsMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
