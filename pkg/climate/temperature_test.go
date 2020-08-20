package climate

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math"
	"testing"
)

func TestEquilibrium(t *testing.T) {
	// Ensure that the equilibrium max temperature of the pole is 0 C.
	got := PoleEquilibrium(OceanSpecificHeat, 23.5 * math.Pi / 180) - ZeroCelsius

	if diff := cmp.Diff(0.0, got, cmpopts.EquateApprox(0.0, 0.1)); diff != "" {
		t.Error(diff)
	}
}

func TestClimate_LowHigh(t *testing.T) {
	tcs := []struct{
		name string
		climate Climate
		latitude float64
		wantLow float64
		wantHigh float64
	} {
		{
			name: "Equatorial Ocean",
			climate: Climate{
				LandSpecificHeat: OceanSpecificHeat,
			},
			latitude: 0,
			wantLow: 27.0,
			wantHigh: 31.0,
		},
		{
			name: "Equatorial Coast",
			climate: Climate{
				LandSpecificHeat: CoastSpecificHeat,
			},
			latitude: 0,
			wantLow: 22.9,
			wantHigh: 35.0,
		},
		{
			name: "Equatorial Desert",
			climate: Climate{
				LandSpecificHeat: DesertSpecificHeat,
			},
			latitude: 0,
			wantLow: 13.8,
			wantHigh: 44.0,
		},
		// Temperate latitudes are better at retaining heat.
		// The atmosphere is more opaque at lower temperatures.
		{
			name: "Temperate Ocean",
			climate: Climate{
				LandSpecificHeat: OceanSpecificHeat,
			},
			latitude: 40,
			wantLow: -4.6,
			wantHigh: -1.5,
		},
		{
			name: "Temperate Coast",
			climate: Climate{
				LandSpecificHeat: CoastSpecificHeat,
			},
			latitude: 40,
			wantLow: -7.7,
			wantHigh: 1.6,
		},
		{
			name: "Temperate Desert",
			climate: Climate{
				LandSpecificHeat: DesertSpecificHeat,
			},
			latitude: 40,
			wantLow: -14.7,
			wantHigh: 8.4,
		},
		{
			name: "Arctic Ocean",
			climate: Climate{
				LandSpecificHeat: OceanSpecificHeat,
			},
			latitude: 70,
			wantLow: -67,
			wantHigh: -65.7,
		},
		{
			name: "Arctic Coast",
			climate: Climate{
				LandSpecificHeat: CoastSpecificHeat,
			},
			latitude: 70,
			wantLow: -68.5,
			wantHigh: -64.3,
		},
		{
			name: "Arctic Desert",
			climate: Climate{
				LandSpecificHeat: DesertSpecificHeat,
			},
			latitude: 70,
			wantLow: -71.6,
			wantHigh: -61.2,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			latitude := tc.latitude * math.Pi / 180.0

			gotLow, gotHigh := LowHigh(tc.climate.LandSpecificHeat, latitude, ZeroCelsius)
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
