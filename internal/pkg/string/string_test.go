package pkg

import (
	"reflect"
	"testing"
)

func TestRemoveSubstring(t *testing.T) {
	type args struct {
		s      string
		substr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should remove sub string from string",
			args: args{
				s:      "hello girl, my name is x",
				substr: "x",
			},
			want: "hello girl, my name is",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveSubstring(tt.args.s, tt.args.substr); got != tt.want {
				t.Errorf("RemoveSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicateValues(t *testing.T) {
	type args struct {
		strSlice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should remove duplicate values",
			args: args{
				strSlice: []string{"123a", "123b", "123c", "123d", "123c", "123b"},
			},
			want: []string{"123a", "123b", "123c", "123d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicateValues(tt.args.strSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicateValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test parse url",
			args: args{
				input: "https://dactoankmapydev.github.io/",
			},
			want:    "dactoankmapydev.github.io",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
