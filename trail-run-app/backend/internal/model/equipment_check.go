package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EquipmentCheck struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:char(36);not null;index:idx_user_race_equip" json:"user_id"`
	RaceID     uuid.UUID `gorm:"type:char(36);not null;index:idx_user_race_equip" json:"race_id"`
	EquipmentID uuid.UUID `gorm:"type:char(36);not null;index:idx_user_race_equip" json:"equipment_id"`
	Checked     bool      `gorm:"default:false" json:"checked"`
}

func (e *EquipmentCheck) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// FavoriteRace - 用户收藏的赛事
type FavoriteRace struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:idx_user_race" json:"user_id"`
	RaceID       uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:idx_user_race" json:"race_id"`
	ReminderSet  bool      `gorm:"default:false" json:"reminder_set"`
	CreatedAt    time.Time `json:"created_at"`
}

func (f *FavoriteRace) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
