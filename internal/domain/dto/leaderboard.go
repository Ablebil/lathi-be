package dto

import "github.com/google/uuid"

type LeaderboardItemResponse struct {
	Rank      int       `json:"rank"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Title     string    `json:"title"`
	Score     int       `json:"score"`
}

type LeaderboardResponse struct {
	TopUsers []LeaderboardItemResponse `json:"top_users"`
}
