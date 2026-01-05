package dto

import (
	"time"

	"github.com/google/uuid"
)

type EditUserProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
}

type UserBadgeResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url"`
	EarnedAt    time.Time `json:"earned_at"`
}

type UserStatsResponse struct {
	TotalChapters     int64   `json:"total_chapters"`
	CompletedChapters int     `json:"completed_chapters"`
	ProgressPercent   float64 `json:"progress_percent"`
	TotalVocabs       int64   `json:"total_vocabs"`
	CollectedVocabs   int     `json:"collected_vocabs"`
}

type UserProfileResponse struct {
	ID           uuid.UUID           `json:"id"`
	Username     string              `json:"username"`
	Email        string              `json:"email"`
	AvatarURL    string              `json:"avatar_url"`
	CurrentTitle string              `json:"current_title"`
	Stats        UserStatsResponse   `json:"stats"`
	Badges       []UserBadgeResponse `json:"badges"`
}
