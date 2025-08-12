package templates

import (
	"testing"
	"time"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/assert"
)

// func TestHumanDate(t *testing.T) {
// 	tm := time.Date(2025, 8, 11, 4, 30, 0, 0, time.UTC)
// 	hd := humanDate(tm)

// 	if hd != "11 Aug 2025 at 04:30" {
// 		t.Errorf("got %q; want %q", hd, "11 Aug 2025 at 04:30")
// 	}
// }

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 8, 11, 4, 30, 0, 0, time.UTC),
			want: "11 Aug 2025 at 04:30",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
