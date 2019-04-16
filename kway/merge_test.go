package kway

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	type args struct {
		arrs [][]int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			"simple",
			args{
				[][]int64{
					{1, 3, 5, 7},
					{2, 4, 6, 8},
				},
			},
			[]int64{1, 2, 3, 4, 5, 6, 7, 8},
		}, {
			"four arrays",
			args{
				[][]int64{
					{1, 5, 9, 13},
					{2, 6, 10, 14},
					{3, 7, 11, 15},
					{4, 8, 12, 16},
				},
			},
			[]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		}, {
			"variant-length arrays",
			args{
				[][]int64{
					{1, 5, 9, 13, 17, 18},
					{2, 6, 10, 14, 19, 20},
					{3, 7, 11, 15},
					{4, 8, 12, 16},
					{21, 22, 23},
				},
			},
			[]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.arrs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type args struct {
		data []int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sort(tt.args.data)
		})
	}
}
