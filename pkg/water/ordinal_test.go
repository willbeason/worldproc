package water

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestOrdinalList_Insert(t *testing.T) {
	tcs := []struct{
		name string
		insert Ordinal
		start OrdinalList
		want OrdinalList
	} {
		{
			name: "insert into empty list",
			insert: Ordinal{Height: 1.0},
			want: OrdinalList{{Height: 1.0}},
		},
		{
			name: "prepend to length 1",
			insert: Ordinal{Height: 0.0},
			start: OrdinalList{{Height: 1.0}},
			want: OrdinalList{{Height: 0.0}, {Height: 1.0}},
		},
		{
			name: "append to length 1",
			insert: Ordinal{Height: 2.0},
			start: OrdinalList{{Height: 1.0}},
			want: OrdinalList{{Height: 1.0}, {Height: 2.0}},
		},
		{
			name: "middle of length 2",
			insert: Ordinal{Height: 1.0},
			start: OrdinalList{{Height: 0.0}, {Height: 2.0}},
			want: OrdinalList{{Height: 0.0}, {Height: 1.0}, {Height: 2.0}},
		},
		{
			name: "middle of length 4",
			insert: Ordinal{Height: 2.0},
			start: OrdinalList{{Height: 0.0}, {Height: 1.0}, {Height: 3.0}, {Height: 4.0}},
			want: OrdinalList{{Height: 0.0}, {Height: 1.0}, {Height: 2.0}, {Height: 3.0}, {Height: 4.0}},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			got.Insert(tc.insert)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
