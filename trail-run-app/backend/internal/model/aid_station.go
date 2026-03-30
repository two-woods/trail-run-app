package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AidStation struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	RaceID     uuid.UUID `gorm:"type:char(36);not null;index" json:"race_id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	DistanceKM float64  `gorm:"type:decimal(6,2);not null" json:"distance_km"`
	ElevationM int       `gorm:"not null" json:"elevation_m"`
	Lat        float64   `gorm:"type:decimal(10,7);not null" json:"lat"`
	Lng        float64   `gorm:"type:decimal(10,7);not null" json:"lng"`
	Supplies   string    `gorm:"type:json;not null" json:"supplies"` // JSON array
	CloseTime  string    `gorm:"type:time;not null" json:"close_time"`
}

func (a *AidStation) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (a *AidStation) GetSupplies() ([]string, error) {
	var supplies []string
	if err := json.Unmarshal([]byte(a.Supplies), &supplies); err != nil {
		return nil, err
	}
	return supplies, nil
}

func (a *AidStation) SetSupplies(supplies []string) error {
	data, err := json.Marshal(supplies)
	if err != nil {
		return err
	}
	a.Supplies = string(data)
	return nil
}

type AidStationResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	DistanceKM float64  `json:"distance_km"`
	ElevationM int       `json:"elevation_m"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Supplies   []string  `json:"supplies"`
	CloseTime  string    `json:"close_time"`
}

func (a *AidStation) ToResponse() (*AidStationResponse, error) {
	supplies, err := a.GetSupplies()
	if err != nil {
		return nil, err
	}
	return &AidStationResponse{
		ID:         a.ID,
		Name:       a.Name,
		DistanceKM: a.DistanceKM,
		ElevationM: a.ElevationM,
		Lat:        a.Lat,
		Lng:        a.Lng,
		Supplies:   supplies,
		CloseTime:  a.CloseTime,
	}, nil
}
