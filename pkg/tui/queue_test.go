package tui

import (
	"reflect"
	"testing"
)

func Test_priorityQueue(t *testing.T) {
	tests := []struct {
		name string
		args []int8
		want []int8
	}{
		{
			name: "PopView returns views in order of priority",
			args: []int8{3, -2, 5, 0},
			want: []int8{5, 3, 0, -2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := newQueue()
			for _, arg := range tt.args {
				q.PushView(&View{priority: arg})
			}
			got := make([]int8, 0)
			for q.Len() > 0 {
				got = append(got, q.PopView().priority)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
