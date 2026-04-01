package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Hotel 酒店（用于缓存赛事周边的酒店信息）
type Hotel struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	Lat             float64   `gorm:"type:decimal(10,7)" json:"lat"`
	Lng             float64   `gorm:"type:decimal(10,7)" json:"lng"`
	DistanceFromStart int      `gorm:"type:int" json:"distance_from_start"` // 距起点距离(米)
	Address         string    `gorm:"type:varchar(500)" json:"address"`
	Phone           string    `gorm:"type:varchar(50)" json:"phone"`
	BusinessArea    string    `gorm:"type:varchar(255)" json:"business_area"` // 商圈
	CityName        string    `gorm:"type:varchar(100)" json:"city_name"`
}

func (h *Hotel) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// FavoriteHotel 用户收藏的酒店
type FavoriteHotel struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index" json:"user_id"`
	RaceID    uuid.UUID `gorm:"type:char(36);not null;index" json:"race_id"` // 关联赛事
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Lat       float64   `gorm:"type:decimal(10,7)" json:"lat"`
	Lng       float64   `gorm:"type:decimal(10,7)" json:"lng"`
	Distance  int       `gorm:"type:int" json:"distance"` // 距起点距离(米)
	Address   string    `gorm:"type:varchar(500)" json:"address"`
	Phone     string    `gorm:"type:varchar(50)" json:"phone"`
	BookingURL string   `gorm:"type:varchar(500)" json:"booking_url"` // 预订链接
	Notes     string    `gorm:"type:text" json:"notes"` // 用户备注
	CreatedAt string    `json:"created_at"`
}

func (f *FavoriteHotel) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// HotelResponse 返回给前端的酒店信息
type HotelResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	DistanceFromStart int    `json:"distance_from_start"`
	Address          string `json:"address"`
	Phone            string `json:"phone,omitempty"`
	BookingURL       string `json:"booking_url,omitempty"`
	Source           string `json:"source"` // "api" 或 "favorite"
}

func (h *Hotel) ToResponse() *HotelResponse {
	return &HotelResponse{
		ID:               h.ID.String(),
		Name:             h.Name,
		Lat:              h.Lat,
		Lng:              h.Lng,
		DistanceFromStart: h.DistanceFromStart,
		Address:          h.Address,
		Phone:            h.Phone,
		Source:           "api",
	}
}

func (f *FavoriteHotel) ToResponse() *HotelResponse {
	return &HotelResponse{
		ID:               f.ID.String(),
		Name:             f.Name,
		Lat:              f.Lat,
		Lng:              f.Lng,
		DistanceFromStart: f.Distance,
		Address:          f.Address,
		Phone:            f.Phone,
		BookingURL:       f.BookingURL,
		Source:           "favorite",
	}
}
