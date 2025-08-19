//go:build js && wasm

package d1

import (
	"math"
	"testing"
)

func Test_isIntegralNumber(t *testing.T) {
	tests := map[string]struct {
		f    float64
		want bool
	}{
		"valid positive integral value": {
			f:    1,
			want: true,
		},
		"valid negative integral value": {
			f:    -1,
			want: true,
		},
		"invalid positive float value": {
			f:    1.1,
			want: false,
		},
		"invalid negative float value": {
			f:    -1.1,
			want: false,
		},
		"invalid NaN": {
			f:    math.NaN(),
			want: false,
		},
		"invalid +Inf": {
			f:    math.Inf(+1),
			want: false,
		},
		"invalid -Inf": {
			f:    math.Inf(-1),
			want: false,
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := isIntegralNumber(tc.f); got != tc.want {
				t.Errorf("isIntegralNumber() = %v, want %v", got, tc.want)
			}
		})
	}
}
