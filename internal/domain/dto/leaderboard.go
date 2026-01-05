package dto

import "github.com/google/uuid"

type LeaderboardItemResponse struct {
	Rank     int       `json:"rank"`
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Title    string    `json:"title"`
	Score    int       `json:"score"`
}

type LeaderboardResponse struct {
	TopUsers []LeaderboardItemResponse `json:"top_users"`
}
