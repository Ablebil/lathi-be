package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Badge struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;not null"`
	Code        string    `json:"code" gorm:"type:varchar(50);unique;not null"` // unique identifier
	Name        string    `json:"name" gorm:"type:varchar(50);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	IconURL     string    `json:"icon_url" gorm:"type:varchar(255);not null"`
}

func (b *Badge) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = id
	}
	return nil
}

type UserBadge struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	BadgeID  uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	EarnedAt time.Time `gorm:"autoCreateTime;not null"`

	User  User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Badge Badge `gorm:"foreignKey:BadgeID;references:ID;constraint:OnDelete:CASCADE"`
}
