package handler

import (
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type leaderboardHandler struct {
	uc contract.LeaderboardUsecaseItf
}

func NewLeaderboardHandler(router fiber.Router, lbUc contract.LeaderboardUsecaseItf) {
	handler := &leaderboardHandler{
		uc: lbUc,
	}

	leaderboardRouter := router.Group("/leaderboards")
	leaderboardRouter.Get("/", handler.getLeaderboard)
}

func (h *leaderboardHandler) getLeaderboard(ctx *fiber.Ctx) error {
	resp, apiErr := h.uc.GetLeaderboard(ctx.Context())
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Leaderboard berhasil dimuat", resp)
}
