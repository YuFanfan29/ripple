package timeline_test

import (
	"slices"
	"testing"

	"github.com/YuFanfan29/ripple/services/feed/internal/timeline"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name    string
		sources [][]int64
		cursor  int64
		limit   int
		want    timeline.Page
	}{
		{
			name:    "one source within the limit",
			sources: [][]int64{{9, 7, 5}},
			limit:   5,
			want:    timeline.Page{PostIDs: []int64{9, 7, 5}},
		},
		{
			name:    "sources are merged newest first",
			sources: [][]int64{{9, 5}, {8, 6}},
			limit:   5,
			want:    timeline.Page{PostIDs: []int64{9, 8, 6, 5}},
		},
		{
			name:    "an ID in two sources appears once",
			sources: [][]int64{{9, 7}, {7}},
			limit:   5,
			want:    timeline.Page{PostIDs: []int64{9, 7}},
		},
		{
			name:    "only IDs older than the cursor",
			sources: [][]int64{{9, 7, 5, 3}},
			cursor:  7,
			limit:   5,
			want:    timeline.Page{PostIDs: []int64{5, 3}},
		},
		{
			name:    "more IDs than the limit",
			sources: [][]int64{{9, 8, 7, 6}},
			limit:   3,
			want:    timeline.Page{PostIDs: []int64{9, 8, 7}, NextCursor: 7},
		},
		{
			name:    "exactly limit IDs is the last page",
			sources: [][]int64{{9, 8, 7}},
			limit:   3,
			want:    timeline.Page{PostIDs: []int64{9, 8, 7}},
		},
		{
			name:    "no sources",
			sources: nil,
			limit:   3,
			want:    timeline.Page{},
		},
		{
			name:    "every source is empty",
			sources: [][]int64{{}, {}},
			limit:   3,
			want:    timeline.Page{},
		},
		{
			// Bob's inbox, a celebrity account he follows, and his own posts.
			name:    "first page of the example",
			sources: [][]int64{{105, 101, 98, 90}, {103, 95}, {101}},
			limit:   4,
			want:    timeline.Page{PostIDs: []int64{105, 103, 101, 98}, NextCursor: 98},
		},
		{
			// Posts 107 and 106 arrived after the first page was read.
			name:    "second page of the example ignores newer posts",
			sources: [][]int64{{107, 105, 101, 98, 90}, {106, 103, 95}, {101}},
			cursor:  98,
			limit:   4,
			want:    timeline.Page{PostIDs: []int64{95, 90}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timeline.Merge(tt.sources, tt.cursor, tt.limit)
			if !slices.Equal(got.PostIDs, tt.want.PostIDs) || got.NextCursor != tt.want.NextCursor {
				t.Errorf("Merge(%v, %d, %d) = %+v, want %+v",
					tt.sources, tt.cursor, tt.limit, got, tt.want)
			}
		})
	}
}

func TestMergePanicsWhenLimitIsLessThanOne(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Merge did not panic")
		}
	}()
	timeline.Merge([][]int64{{9}}, 0, 0)
}
