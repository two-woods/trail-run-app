package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EquipmentCategory string

const (
	EquipmentCategory穿着 EquipmentCategory = "穿着"
	EquipmentCategory装备 EquipmentCategory = "装备"
	EquipmentCategory补给 EquipmentCategory = "补给"
	EquipmentCategory安全 EquipmentCategory = "安全"
	EquipmentCategory其他 EquipmentCategory = "其他"
)

type Equipment struct {
	ID          uuid.UUID         `gorm:"type:char(36);primaryKey" json:"id"`
	RaceID     uuid.UUID         `gorm:"type:char(36);not null;index" json:"race_id"`
	Name       string            `gorm:"type:varchar(100);not null" json:"name"`
	IsMandatory bool             `gorm:"default:false" json:"is_mandatory"`
	Category   EquipmentCategory `gorm:"type:varchar(20);not null" json:"category"`
}

func (e *Equipment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type EquipmentResponse struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	IsMandatory bool              `json:"is_mandatory"`
	Category    EquipmentCategory `json:"category"`
	Checked     bool              `json:"checked,omitempty"`
}

func (e *Equipment) ToResponse() *EquipmentResponse {
	return &EquipmentResponse{
		ID:          e.ID,
		Name:        e.Name,
		IsMandatory: e.IsMandatory,
		Category:    e.Category,
	}
}
