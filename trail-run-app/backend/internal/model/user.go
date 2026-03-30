package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Nickname     string    `gorm:"type:varchar(50);not null" json:"nickname"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL    *string   `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	Phone        *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	ITRAAccount  *string   `gorm:"type:varchar(100)" json:"itra_account,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Nickname    string    `json:"nickname"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	ITRAAccount *string   `json:"itra_account,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Nickname:    u.Nickname,
		AvatarURL:   u.AvatarURL,
		Phone:       u.Phone,
		ITRAAccount: u.ITRAAccount,
		CreatedAt:   u.CreatedAt,
	}
}
