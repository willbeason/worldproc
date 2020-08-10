package sun

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"testing"
)

func TestDirectional_DeclinationAzimuth(t *testing.T) {
	tcs := []struct{
		name string
		date float64
		a geodesic.Angle
		want geodesic.Angle
	} {
		{
			name: "just north of equator",
			a: geodesic.Angle{
				Theta: 0.01,
				Phi: 0,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   -math.Pi/2, // points south
			},
		},
		{
			name: "just north east of equator",
			a: geodesic.Angle{
				Theta: 0.01,
				Phi: -0.01,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   -math.Pi/4, // points south-west
			},
		},
		{
			name: "just east of equator",
			a: geodesic.Angle{
				Theta: 0,
				Phi: -0.01,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   0, // points west
			},
		},
		{
			name: "just south east of equator",
			a: geodesic.Angle{
				Theta: -0.01,
				Phi: -0.01,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   math.Pi / 4, // points north-west
			},
		},
		{
			name: "just south of equator",
			a: geodesic.Angle{
				Theta: -0.01,
				Phi: 0,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   math.Pi/2, // points north
			},
		},
		{
			name: "just south west of equator",
			a: geodesic.Angle{
				Theta: -0.01,
				Phi: 0.01,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   3 * math.Pi/4, // points north east
			},
		},
		{
			name: "just west of equator",
			a: geodesic.Angle{
				Theta: 0,
				Phi: 0.01,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   math.Pi, // points east
			},
		},
		{
			name: "just north west of equator",
			a: geodesic.Angle{
				Theta: 0.01,
				Phi: 0.01,
			},
			want: geodesic.Angle{
				Theta: math.Pi / 2,
				Phi:   5 * math.Pi/4, // points south-east
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s := Directional{}
			s.Set(tc.date)

			got := s.AltitudeAzimuth(tc.a)
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.0, 0.02)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
