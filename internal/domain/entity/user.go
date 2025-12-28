package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Title string

const (
	Cantrik Title = "cantrik"
	Abdi    Title = "abdi"
	Priyayi Title = "priyayi"
)

type User struct {
	ID                   uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;not null"`
	Username             string    `json:"username" gorm:"type:varchar(50);unique;not null"`
	Email                string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password             string    `json:"password" gorm:"type:varchar(255);not null"`
	CurrentTitle         Title     `json:"current_title" gorm:"type:varchar(255);default:'cantrik';not null"`
	LastChapterCompleted int       `json:"last_chapter_completed" gorm:"type:int;default:0;not null"`
	TotalWordsCollected  int       `json:"total_words_collected" gorm:"type:int;default:0;not null"`
	IsVerified           bool      `json:"is_verified" gorm:"type:boolean;default:false;not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:timestamp;autoCreateTime;not null"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"type:timestamp;autoUpdateTime;not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}
