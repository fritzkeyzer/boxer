package boxer

import (
	"reflect"
	"testing"
)

func ExampleBox_Print() {
	b := New()
	b.WrapLimit = 50
	b.AddHeader("id", "name", "value")
	b.AddLine(1, "foo", "bar")
	b.AddLine(2, "example", "multi\nline")
	b.Print()

	//output:id   name      value
	//1    foo       bar
	//2    example   multi
	//               line
}

func ExampleBox_PrintWithBorders() {
	b := NewWithBorders()
	b.WrapLimit = 50
	b.AddHeader("id", "name", "value")
	b.AddLine(1, "foo", "bar")
	b.AddLine(2, "example", "multi\nline")
	b.Print()

	//output:┌──┬───────┬─────┐
	//│id│name   │value│
	//├──┼───────┼─────┤
	//│1 │foo    │bar  │
	//│2 │example│multi│
	//│  │       │line │
	//└──┴───────┴─────┘
}

func Test_wrapLine(t *testing.T) {
	type args struct {
		in    string
		limit int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			"A",
			args{
				in:    "0123456789",
				limit: 3,
			},
			[]string{"012", "345", "678", "9"},
		},
		{
			"B",
			args{
				in:    "0123456789",
				limit: 10,
			},
			[]string{"0123456789"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrapLine(tt.args.in, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
