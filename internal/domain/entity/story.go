package entity

import (
	"time"

	"github.com/Ablebil/lathi-be/internal/domain/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chapter struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;not null"`
	Title         string    `json:"title" gorm:"type:varchar(100);not null"`
	Description   string    `json:"description" gorm:"type:text;not null"`
	CoverImageURL string    `json:"cover_image_url" gorm:"type:varchar(255);not null"`
	OrderIndex    int       `json:"order_index" gorm:"type:int;not null"`

	Slides []Slide `json:"slides" gorm:"foreignKey:ChapterID;references:ID;constraint:OnDelete:CASCADE"`
}

func (c *Chapter) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		c.ID = id
	}
	return nil
}

type Slide struct {
	ID                 uuid.UUID   `json:"id" gorm:"type:char(36);primaryKey;not null"`
	ChapterID          uuid.UUID   `json:"chapter_id" gorm:"type:char(36);not null"`
	BackgroundImageURL string      `json:"background_image_url" gorm:"type:varchar(255);not null"`
	Characters         types.JSONB `json:"characters" gorm:"type:jsonb;default:'[]'::jsonb;not null"`
	SpeakerName        string      `json:"speaker_name" gorm:"type:varchar(100);not null"`
	Content            string      `json:"content" gorm:"type:text;not null"`
	NextSlideID        *uuid.UUID  `json:"next_slide_id" gorm:"type:char(36)"`
	Choices            types.JSONB `json:"choices" gorm:"type:jsonb;default:'[]'::jsonb"`

	Vocabularies []Dictionary `json:"vocabularies" gorm:"many2many:slide_vocabularies;constraint:OnDelete:CASCADE"`
}

func (s *Slide) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		s.ID = id
	}
	return nil
}

type UserStorySession struct {
	ID             uuid.UUID   `json:"id" gorm:"type:char(36);primaryKey;not null"`
	UserID         uuid.UUID   `json:"user_id" gorm:"type:char(36);not null;uniqueIndex:idx_user_chapter"`
	ChapterID      uuid.UUID   `json:"chapter_id" gorm:"type:char(36);not null;uniqueIndex:idx_user_chapter"`
	CurrentSlideID uuid.UUID   `json:"current_slide_id" gorm:"type:char(36);not null"`
	CurrentHearts  int         `json:"current_hearts" gorm:"type:int;default:3;not null"`
	IsGameOver     bool        `json:"is_game_over" gorm:"type:boolean;default:false;not null"`
	IsCompleted    bool        `json:"is_completed" gorm:"type:boolean;default:false;not null"`
	HistoryLog     types.JSONB `json:"history_log" gorm:"type:jsonb;default:'[]'::jsonb;not null"`
	CreatedAt      time.Time   `json:"created_at" gorm:"type:timestamp;autoCreateTime;not null"`
	UpdatedAt      time.Time   `json:"updated_at" gorm:"type:timestamp;autoUpdateTime;not null"`

	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Chapter Chapter `gorm:"foreignKey:ChapterID;references:ID;constraint:OnDelete:CASCADE"`
}

func (uss *UserStorySession) BeforeCreate(tx *gorm.DB) error {
	if uss.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		uss.ID = id
	}
	return nil
}
