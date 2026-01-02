package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dictionary struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;not null"`
	WordKrama string    `json:"word_krama" gorm:"type:varchar(100);not null"`
	WordNgoko string    `json:"word_ngoko" gorm:"type:varchar(100);not null"`
	WordIndo  string    `json:"word_indo" gorm:"type:varchar(100);not null"`
}

func (d *Dictionary) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		d.ID = id
	}
	return nil
}

type UserVocabulary struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	DictionaryID uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	UnlockedAt   time.Time `gorm:"autoCreateTime;not null"`

	User       User       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Dictionary Dictionary `gorm:"foreignKey:DictionaryID;references:ID;constraint:OnDelete:CASCADE"`
}
