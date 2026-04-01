package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GetMockDataPath returns the absolute path to the testdata directory
func GetMockDataPath() string {
	// Get the directory where this file is located
	dir, _ := filepath.Split(os.Args[0])
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, "internal", "testdata")
}

// LoadMockRaces loads mock races data
func LoadMockRaces() (*MockRacesResponse, error) {
	path := filepath.Join(GetMockDataPath(), "mock_races.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var response MockRacesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// LoadMockHotels loads mock hotels data
func LoadMockHotels() (*MockHotelsResponse, error) {
	path := filepath.Join(GetMockDataPath(), "mock_hotels.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var response MockHotelsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// LoadMockResults loads mock results data
func LoadMockResults() (*MockResultsResponse, error) {
	path := filepath.Join(GetMockDataPath(), "mock_results.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var response MockResultsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Mock data structures matching the JSON files

type MockRacesResponse struct {
	Races   []MockRace `json:"races"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	PageSize int       `json:"page_size"`
}

type MockRace struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Date            string  `json:"date"`
	Location        string  `json:"location"`
	Province        string  `json:"province"`
	City            string  `json:"city"`
	DistanceKM      float64 `json:"distance_km"`
	ElevationM      int     `json:"elevation_m"`
	Difficulty      string  `json:"difficulty"`
	ITRAPoints     *float64 `json:"itra_points,omitempty"`
	StartLat       float64 `json:"start_lat"`
	StartLng       float64 `json:"start_lng"`
	EndLat         float64 `json:"end_lat"`
	EndLng         float64 `json:"end_lng"`
	RouteGPXURL    *string `json:"route_gpx_url,omitempty"`
	WeatherCityCode *string `json:"weather_city_code,omitempty"`
	Status         string  `json:"status"`
}

type MockHotelsResponse struct {
	APIHotels      []MockHotel `json:"api_hotels"`
	FavoriteHotels []MockHotel `json:"favorite_hotels"`
	TotalAPI      int         `json:"total_api"`
	TotalFavorite int         `json:"total_favorite"`
}

type MockHotel struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	DistanceFromStart int    `json:"distance_from_start"`
	Address          string  `json:"address"`
	Phone            string  `json:"phone,omitempty"`
	BusinessArea     string  `json:"business_area,omitempty"`
	CityName         string  `json:"city_name,omitempty"`
	BookingURL       string  `json:"booking_url,omitempty"`
	Notes            string  `json:"notes,omitempty"`
	Source           string  `json:"source"` // "api" or "favorite"
}

type MockResultsResponse struct {
	Results  []MockResult `json:"results"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type MockResult struct {
	ID                 string   `json:"id"`
	UserID             string   `json:"user_id"`
	RaceID             string   `json:"race_id"`
	RaceName           string   `json:"race_name"`
	RaceDate           string   `json:"race_date"`
	DistanceKM         float64  `json:"distance_km"`
	ElevationM         int      `json:"elevation_m"`
	FinishTime         *string  `json:"finish_time,omitempty"`
	Ranking            *int     `json:"ranking,omitempty"`
	RankingAgeGroup    *int     `json:"ranking_age_group,omitempty"`
	GPXURL             *string  `json:"gpx_url,omitempty"`
	Photos             []string `json:"photos"`
	GeneratedImageURL  *string  `json:"generated_image_url,omitempty"`
	Status             string   `json:"status"`
}
