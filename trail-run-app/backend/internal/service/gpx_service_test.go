package service

import (
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		minDist  float64
		maxDist  float64
	}{
		{
			name:    "Same point",
			lat1:    30.5728,
			lon1:    114.2525,
			lat2:    30.5728,
			lon2:    114.2525,
			minDist: 0,
			maxDist: 0.001,
		},
		{
			name:    "Wuhan to Beijing approx",
			lat1:    30.5728,
			lon1:    114.2525,
			lat2:    39.9042,
			lon2:    116.4074,
			minDist: 1050,
			maxDist: 1100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist := haversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if dist < tt.minDist || dist > tt.maxDist {
				t.Errorf("haversineDistance() = %v, want between %v and %v", dist, tt.minDist, tt.maxDist)
			}
		})
	}
}

func TestParseGPX(t *testing.T) {
	// Valid GPX with one track point
	validGPX := `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="Test">
  <trk>
    <name>Test Track</name>
    <trkseg>
      <trkpt lat="30.5728" lon="114.2525">
        <ele>100</ele>
      </trkpt>
      <trkpt lat="30.5730" lon="114.2528">
        <ele>110</ele>
      </trkpt>
      <trkpt lat="30.5732" lon="114.2530">
        <ele>105</ele>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`

	stats, geojson, err := ParseGPX(validGPX)
	if err != nil {
		t.Fatalf("ParseGPX() error = %v", err)
	}

	if stats == nil {
		t.Fatal("ParseGPX() returned nil stats")
	}

	if geojson == nil {
		t.Fatal("ParseGPX() returned nil geojson")
	}

	if stats.DistanceKM <= 0 {
		t.Errorf("ParseGPX() distance = %v, want > 0", stats.DistanceKM)
	}

	if len(geojson.Features) == 0 {
		t.Error("ParseGPX() returned no features")
	}
}

func TestParseGPXInvalid(t *testing.T) {
	invalidGPX := `<?xml version="1.0"?>
<invalid>gpx content</invalid>`

	_, _, err := ParseGPX(invalidGPX)
	if err == nil {
		t.Error("ParseGPX() expected error for invalid GPX")
	}
}

func TestParseGPXEmpty(t *testing.T) {
	emptyGPX := `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="Test">
</gpx>`

	stats, geojson, err := ParseGPX(emptyGPX)
	if err != nil {
		t.Fatalf("ParseGPX() unexpected error = %v", err)
	}

	if stats == nil {
		t.Fatal("ParseGPX() returned nil stats")
	}

	// Empty track should still return valid but empty geojson
	if geojson == nil {
		t.Fatal("ParseGPX() returned nil geojson")
	}
}

func TestGPXStatsElevation(t *testing.T) {
	// GPX with elevation changes
	gpxWithElevation := `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="Test">
  <trk>
    <name>Elevation Test</name>
    <trkseg>
      <trkpt lat="30.5728" lon="114.2525">
        <ele>100</ele>
      </trkpt>
      <trkpt lat="30.5730" lon="114.2528">
        <ele>150</ele>
      </trkpt>
      <trkpt lat="30.5732" lon="114.2530">
        <ele>120</ele>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`

	stats, _, err := ParseGPX(gpxWithElevation)
	if err != nil {
		t.Fatalf("ParseGPX() error = %v", err)
	}

	// Elevation gain should be 50m (100->150)
	if stats.ElevationGainM < 40 || stats.ElevationGainM > 60 {
		t.Errorf("ElevationGainM = %v, want ~50", stats.ElevationGainM)
	}

	// Elevation loss should be 30m (150->120)
	if stats.ElevationLossM < 20 || stats.ElevationLossM > 40 {
		t.Errorf("ElevationLossM = %v, want ~30", stats.ElevationLossM)
	}

	// Max elevation should be 150
	if stats.MaxElevationM != 150 {
		t.Errorf("MaxElevationM = %v, want 150", stats.MaxElevationM)
	}

	// Min elevation should be 100
	if stats.MinElevationM != 100 {
		t.Errorf("MinElevationM = %v, want 100", stats.MinElevationM)
	}
}
