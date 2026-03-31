package service

import (
	"io"
	"math"
	"net/http"

	"github.com/tkrajina/gpxgo/gpx"
)

// GPXStats holds computed statistics from a GPX track
type GPXStats struct {
	DistanceKM     float64 `json:"distance_km"`      // Total distance in km
	ElevationGainM int     `json:"elevation_gain_m"` // Total ascent in m
	ElevationLossM int     `json:"elevation_loss_m"` // Total descent in m
	MinElevationM  int     `json:"min_elevation_m"`  // Min altitude in m
	MaxElevationM  int     `json:"max_elevation_m"`  // Max altitude in m
	DurationSecs   int     `json:"duration_secs"`    // Total time in seconds
	AvgSpeedKMH    float64 `json:"avg_speed_kmh"`    // Average speed in km/h
}

// GeoJSONFeature represents a GeoJSON Feature
type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{}  `json:"properties,omitempty"`
	Geometry   map[string]interface{}  `json:"geometry"`
}

// GeoJSONFeatureCollection represents a GeoJSON FeatureCollection
type GeoJSONFeatureCollection struct {
	Type     string            `json:"type"`
	Features []GeoJSONFeature  `json:"features"`
}

// ParseGPX parses GPX data and returns stats and GeoJSON
func ParseGPX(gpxData string) (*GPXStats, *GeoJSONFeatureCollection, error) {
	parsed, err := gpx.ParseString(gpxData)
	if err != nil {
		return nil, nil, err
	}

	stats := &GPXStats{}
	var coordinates [][]float64

	for _, track := range parsed.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				// Haversine distance calculation
				if len(coordinates) > 0 {
					last := coordinates[len(coordinates)-1]
					dist := haversineDistance(last[0], last[1], point.Latitude, point.Longitude)
					stats.DistanceKM += dist
				}

				coordinates = append(coordinates, []float64{point.Longitude, point.Latitude})

				if point.Elevation.NotNull() {
					elevation := int(point.Elevation.Value())
					if stats.MinElevationM == 0 || elevation < stats.MinElevationM {
						stats.MinElevationM = elevation
					}
					if elevation > stats.MaxElevationM {
						stats.MaxElevationM = elevation
					}
				}
			}

			// Calculate elevation gain/loss within segment
			for i := 1; i < len(segment.Points); i++ {
				prevEle := segment.Points[i-1].Elevation
				currEle := segment.Points[i].Elevation
				if prevEle.NotNull() && currEle.NotNull() {
					diff := int(currEle.Value()) - int(prevEle.Value())
					if diff > 0 {
						stats.ElevationGainM += diff
					} else {
						stats.ElevationLossM += int(math.Abs(float64(diff)))
					}
				}
			}
		}
	}

	// Calculate duration if time info available
	if len(parsed.Tracks) > 0 && len(parsed.Tracks[0].Segments) > 0 {
		segment := parsed.Tracks[0].Segments[0]
		if len(segment.Points) >= 2 {
			start := segment.Points[0]
			end := segment.Points[len(segment.Points)-1]
			if start.Timestamp.After(end.Timestamp) == false {
				duration := end.Timestamp.Sub(start.Timestamp)
				stats.DurationSecs = int(duration.Seconds())
				if stats.DurationSecs > 0 {
					stats.AvgSpeedKMH = (stats.DistanceKM / float64(stats.DurationSecs)) * 3600
				}
			}
		}
	}

	// Build GeoJSON LineString
	geojson := &GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]GeoJSONFeature, 0),
	}

	if len(coordinates) > 0 {
		lineString := map[string]interface{}{
			"type":        "LineString",
			"coordinates": coordinates,
		}

		geojson.Features = append(geojson.Features, GeoJSONFeature{
			Type:       "Feature",
			Properties: map[string]interface{}{},
			Geometry:   lineString,
		})

		// Add start point
		if len(coordinates) > 0 {
			geojson.Features = append(geojson.Features, GeoJSONFeature{
				Type: "Feature",
				Properties: map[string]interface{}{
					"marker": "start",
				},
				Geometry: map[string]interface{}{
					"type":        "Point",
					"coordinates": coordinates[0],
				},
			})

			// Add end point
			geojson.Features = append(geojson.Features, GeoJSONFeature{
				Type: "Feature",
				Properties: map[string]interface{}{
					"marker": "end",
				},
				Geometry: map[string]interface{}{
					"type":        "Point",
					"coordinates": coordinates[len(coordinates)-1],
				},
			})
		}
	}

	return stats, geojson, nil
}

// ParseGPXFromURL fetches GPX from a URL and parses it
func ParseGPXFromURL(gpxURL string) (*GPXStats, *GeoJSONFeatureCollection, error) {
	resp, err := http.Get(gpxURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	parsed, err := gpx.ParseBytes(data)
	if err != nil {
		return nil, nil, err
	}

	stats := &GPXStats{}
	var coordinates [][]float64

	for _, track := range parsed.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				if len(coordinates) > 0 {
					last := coordinates[len(coordinates)-1]
					dist := haversineDistance(last[0], last[1], point.Latitude, point.Longitude)
					stats.DistanceKM += dist
				}

				coordinates = append(coordinates, []float64{point.Longitude, point.Latitude})

				if point.Elevation.NotNull() {
					elevation := int(point.Elevation.Value())
					if stats.MinElevationM == 0 || elevation < stats.MinElevationM {
						stats.MinElevationM = elevation
					}
					if elevation > stats.MaxElevationM {
						stats.MaxElevationM = elevation
					}
				}
			}

			for i := 1; i < len(segment.Points); i++ {
				prevEle := segment.Points[i-1].Elevation
				currEle := segment.Points[i].Elevation
				if prevEle.NotNull() && currEle.NotNull() {
					diff := int(currEle.Value()) - int(prevEle.Value())
					if diff > 0 {
						stats.ElevationGainM += diff
					} else {
						stats.ElevationLossM += int(math.Abs(float64(diff)))
					}
				}
			}
		}
	}

	if len(parsed.Tracks) > 0 && len(parsed.Tracks[0].Segments) > 0 {
		segment := parsed.Tracks[0].Segments[0]
		if len(segment.Points) >= 2 {
			start := segment.Points[0]
			end := segment.Points[len(segment.Points)-1]
			if !start.Timestamp.After(end.Timestamp) {
				duration := end.Timestamp.Sub(start.Timestamp)
				stats.DurationSecs = int(duration.Seconds())
				if stats.DurationSecs > 0 {
					stats.AvgSpeedKMH = (stats.DistanceKM / float64(stats.DurationSecs)) * 3600
				}
			}
		}
	}

	geojson := &GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]GeoJSONFeature, 0),
	}

	if len(coordinates) > 0 {
		lineString := map[string]interface{}{
			"type":        "LineString",
			"coordinates": coordinates,
		}

		geojson.Features = append(geojson.Features, GeoJSONFeature{
			Type:       "Feature",
			Properties: map[string]interface{}{},
			Geometry:   lineString,
		})

		geojson.Features = append(geojson.Features, GeoJSONFeature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"marker": "start",
			},
			Geometry: map[string]interface{}{
				"type":        "Point",
				"coordinates": coordinates[0],
			},
		})

		geojson.Features = append(geojson.Features, GeoJSONFeature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"marker": "end",
			},
			Geometry: map[string]interface{}{
				"type":        "Point",
				"coordinates": coordinates[len(coordinates)-1],
			},
		})
	}

	return stats, geojson, nil
}

// haversineDistance calculates distance between two points in km
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}
