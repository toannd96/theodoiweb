package pkg

import (
	"reflect"
	"testing"
)

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
