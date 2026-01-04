package dto

import (
	"time"

	"github.com/google/uuid"
)

type StoryActionRequest struct {
	ChapterID   uuid.UUID `json:"chapter_id" validate:"required,uuid"`
	SlideID     uuid.UUID `json:"slide_id" validate:"required,uuid"`
	ChoiceIndex *int      `json:"choice_index,omitempty"`
}

type HistoryEntry struct {
	Speaker   string    `json:"speaker"`
	Text      string    `json:"text"`
	IsUser    bool      `json:"is_user"`
	Timestamp time.Time `json:"timestamp"`
}

type CharacterOnScreen struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	IsActive bool   `json:"is_active"`
}

type ChapterListReponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CoverImageURL string    `json:"cover_image_url"`
	OrderIndex    int       `json:"order_index"`
	IsLocked      bool      `json:"is_locked"`
	IsCompleted   bool      `json:"is_completed"`
}

type ChapterContentResponse struct {
	ChapterID uuid.UUID           `json:"chapter_id"`
	Slides    []SlideItemResponse `json:"slides"`
}

type SlideItemResponse struct {
	ID                 uuid.UUID            `json:"id"`
	BackgroundImageURL string               `json:"background_image_url"`
	Characters         []CharacterOnScreen  `json:"characters"`
	SpeakerName        string               `json:"speaker_name"`
	Content            string               `json:"content"`
	NextSlideID        *uuid.UUID           `json:"next_slide_id"`
	Vocabularies       []VocabItemResponse  `json:"vocabularies"`
	Choices            []ChoiceItemResponse `json:"choices"`
}

type VocabItemResponse struct {
	ID        uuid.UUID `json:"id"`
	WordKrama string    `json:"word_krama"`
	WordNgoko string    `json:"word_ngoko"`
	WordIndo  string    `json:"word_indo"`
}

type ChoiceItemResponse struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

type UserSessionResponse struct {
	SessionID      uuid.UUID      `json:"session_id"`
	CurrentSlideID uuid.UUID      `json:"current_slide_id"`
	CurrentHearts  int            `json:"current_hearts"`
	IsGameOver     bool           `json:"is_game_over"`
	IsCompleted    bool           `json:"is_completed"`
	HistoryLog     []HistoryEntry `json:"history_log"`
}

type StoryActionResponse struct {
	IsGameOver      bool           `json:"is_game_over"`
	IsCompleted     bool           `json:"is_completed"`
	Message         string         `json:"message"` // msg if gameover/completed
	RemainingHearts int            `json:"remaining_hearts"`
	NextSlideID     *uuid.UUID     `json:"next_slide_id"`
	HistoryLog      []HistoryEntry `json:"history_log"`
}
