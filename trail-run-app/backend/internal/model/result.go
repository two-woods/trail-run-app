package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResultStatus string

const (
	ResultStatusRegistered ResultStatus = "registered"
	ResultStatusOngoing   ResultStatus = "ongoing"
	ResultStatusFinished  ResultStatus = "finished"
)

type Result struct {
	ID                uuid.UUID    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID            uuid.UUID    `gorm:"type:char(36);not null;index" json:"user_id"`
	RaceID            uuid.UUID    `gorm:"type:char(36);not null;index" json:"race_id"`
	FinishTime        *string      `gorm:"type:time" json:"finish_time,omitempty"`
	Ranking           *int         `json:"ranking,omitempty"`
	RankingAgeGroup   *int         `json:"ranking_age_group,omitempty"`
	GPXURL            *string      `gorm:"type:varchar(500)" json:"gpx_url,omitempty"`
	Photos            *string      `gorm:"type:json" json:"photos,omitempty"` // JSON array
	GeneratedImageURL *string      `gorm:"type:varchar(500)" json:"generated_image_url,omitempty"`
	Status            ResultStatus `gorm:"type:varchar(20);default:'registered'" json:"status"`
	CreatedAt         time.Time    `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Race *Race `gorm:"foreignKey:RaceID" json:"race,omitempty"`
}

func (r *Result) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *Result) GetPhotos() ([]string, error) {
	if r.Photos == nil || *r.Photos == "" {
		return []string{}, nil
	}
	var photos []string
	if err := json.Unmarshal([]byte(*r.Photos), &photos); err != nil {
		return nil, err
	}
	return photos, nil
}

func (r *Result) SetPhotos(photos []string) error {
	data, err := json.Marshal(photos)
	if err != nil {
		return err
	}
	r.Photos = new(string)
	*r.Photos = string(data)
	return nil
}

type ResultListItem struct {
	ID                uuid.UUID    `json:"id"`
	RaceName          string       `json:"race_name"`
	RaceDate          time.Time    `json:"race_date"`
	DistanceKM        float64      `json:"distance_km"`
	ElevationM        int          `json:"elevation_m"`
	FinishTime        *string      `json:"finish_time,omitempty"`
	Ranking           *int         `json:"ranking,omitempty"`
	RankingAgeGroup   *int         `json:"ranking_age_group,omitempty"`
	Status            ResultStatus `json:"status"`
	GeneratedImageURL *string      `json:"generated_image_url,omitempty"`
}

type ResultDetail struct {
	ID                uuid.UUID    `json:"id"`
	Race              *Race        `json:"race"`
	FinishTime        *string      `json:"finish_time,omitempty"`
	Ranking           *int         `json:"ranking,omitempty"`
	RankingAgeGroup   *int         `json:"ranking_age_group,omitempty"`
	GPXURL            *string      `json:"gpx_url,omitempty"`
	GeoJSON           interface{}  `json:"geojson,omitempty"`
	Stats             interface{}  `json:"stats,omitempty"`
	Photos            []string     `json:"photos"`
	GeneratedImageURL *string      `json:"generated_image_url,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
}
