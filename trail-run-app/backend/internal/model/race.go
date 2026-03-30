package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "Easy"
	DifficultyMedium Difficulty = "Medium"
	DifficultyHard   Difficulty = "Hard"
	DifficultyExtrem Difficulty = "Extrem"
)

type RaceStatus string

const (
	RaceStatusDraft     RaceStatus = "draft"
	RaceStatusPublished RaceStatus = "published"
)

type Race struct {
	ID              uuid.UUID    `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string       `gorm:"type:varchar(255);not null" json:"name"`
	Date            time.Time    `gorm:"type:date;not null" json:"date"`
	Location        string       `gorm:"type:varchar(255);not null" json:"location"`
	Province        string       `gorm:"type:varchar(50);not null" json:"province"`
	City            string       `gorm:"type:varchar(50);not null" json:"city"`
	DistanceKM      float64      `gorm:"type:decimal(6,2);not null" json:"distance_km"`
	ElevationM      int          `gorm:"not null" json:"elevation_m"`
	Difficulty      Difficulty    `gorm:"type:varchar(10);not null" json:"difficulty"`
	ITRAPoints      *float64     `gorm:"type:decimal(4,1)" json:"itra_points,omitempty"`
	StartLat        float64      `gorm:"type:decimal(10,7);not null" json:"start_lat"`
	StartLng        float64      `gorm:"type:decimal(10,7);not null" json:"start_lng"`
	EndLat          float64      `gorm:"type:decimal(10,7);not null" json:"end_lat"`
	EndLng          float64      `gorm:"type:decimal(10,7);not null" json:"end_lng"`
	RouteGPXURL     *string      `gorm:"type:varchar(500)" json:"route_gpx_url,omitempty"`
	WeatherCityCode *string      `gorm:"type:varchar(20)" json:"weather_city_code,omitempty"`
	Status          RaceStatus   `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CreatedAt       time.Time    `json:"created_at"`

	AidStations []AidStation `gorm:"foreignKey:RaceID" json:"aid_stations,omitempty"`
	Equipment   []Equipment  `gorm:"foreignKey:RaceID" json:"equipment,omitempty"`
}

func (r *Race) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type RaceListItem struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Date         time.Time  `json:"date"`
	Location     string     `json:"location"`
	DistanceKM   float64    `json:"distance_km"`
	ElevationM   int        `json:"elevation_m"`
	Difficulty   Difficulty `json:"difficulty"`
	ITRAPoints   *float64  `json:"itra_points,omitempty"`
	Status       RaceStatus `json:"status"`
}

func (r *Race) ToListItem() *RaceListItem {
	return &RaceListItem{
		ID:         r.ID,
		Name:       r.Name,
		Date:       r.Date,
		Location:   r.Location,
		DistanceKM: r.DistanceKM,
		ElevationM: r.ElevationM,
		Difficulty: r.Difficulty,
		ITRAPoints: r.ITRAPoints,
		Status:     r.Status,
	}
}
