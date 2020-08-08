package climate

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math"
	"testing"
)

func TestClimate_LowHigh(t *testing.T) {
	tcs := []struct{
		name string
		climate Climate
		wantLow float64
		wantHigh float64
	} {
		{
			name: "Equatorial Ocean",
			climate: Climate{
				SpecificHeat: 2000000,
				Latitude:     0,
			},
			wantLow: 27.3,
			wantHigh: 30.4,
		},
		{
			name: "Equatorial Coast",
			climate: Climate{
				SpecificHeat: 1000000,
				Latitude:     0,
			},
			wantLow: 25.8,
			wantHigh: 31.9,
		},
		{
			name: "Equatorial Desert",
			climate: Climate{
				SpecificHeat: 200000,
				Latitude:     0,
			},
			wantLow: 13.8,
			wantHigh: 44.0,
		},
		// Temperate latitudes are better at retaining heat.
		// The atmosphere is more opaque at lower temperatures.
		{
			name: "Temperate Ocean",
			climate: Climate{
				SpecificHeat: 2000000,
				Latitude:     40,
			},
			wantLow: 18.8,
			wantHigh: 21.1,
		},
		{
			name: "Temperate Coast",
			climate: Climate{
				SpecificHeat: 1000000,
				Latitude:     40,
			},
			wantLow: 17.6,
			wantHigh: 22.3,
		},
		{
			name: "Temperate Desert",
			climate: Climate{
				SpecificHeat: 200000,
				Latitude:     40,
			},
			wantLow: 8.4,
			wantHigh: 31.7,
		},
		{
			name: "Arctic Ocean",
			climate: Climate{
				SpecificHeat: 2000000,
				Latitude:     70,
			},
			wantLow: -1.2,
			wantHigh: -0.1,
		},
		{
			name: "Arctic Coast",
			climate: Climate{
				SpecificHeat: 1000000,
				Latitude:     70,
			},
			wantLow: -1.7,
			wantHigh: 0.4,
		},
		{
			name: "Arctic Desert",
			climate: Climate{
				SpecificHeat: 200000,
				Latitude:     70,
			},
			wantLow: -5.8,
			wantHigh: 4.7,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.climate.Latitude *= math.Pi / 180.0

			gotLow, gotHigh := tc.climate.LowHigh(DefaultTemperature)
			gotLow -= 273
			gotHigh -= 273

			if diff := cmp.Diff(tc.wantLow, gotLow, cmpopts.EquateApprox(0.0, 0.1)); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(tc.wantHigh, gotHigh, cmpopts.EquateApprox(0.0, 0.1)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
