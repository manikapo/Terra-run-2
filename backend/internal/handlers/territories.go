package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/territory-run/api/internal/middleware"
	"github.com/territory-run/api/internal/services"
)

type TerritoryHandler struct {
	territories *services.TerritoryService
}

func NewTerritoryHandler(territories *services.TerritoryService) *TerritoryHandler {
	return &TerritoryHandler{territories: territories}
}

func (h *TerritoryHandler) Tile(c *fiber.Ctx) error {
	tileHex := c.Params("h3_tile")
	if tileHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "h3_tile required"})
	}

	geo, err := h.territories.TileGeoJSON(c.Context(), tileHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(geo)
}

func (h *TerritoryHandler) GlobalLeaderboard(c *fiber.Ctx) error {
	board, err := h.territories.GlobalLeaderboard(c.Context(), c.QueryInt("limit", 50))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(board)
}

func (h *TerritoryHandler) LocalLeaderboard(c *fiber.Ctx) error {
	parent := c.Query("h3_parent")
	if parent == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "h3_parent required"})
	}
	board, err := h.territories.LocalLeaderboard(c.Context(), parent, c.QueryInt("limit", 50))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(board)
}

func (h *TerritoryHandler) Events(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	events, err := h.territories.RecentEvents(c.Context(), userID, c.QueryInt("limit", 20))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"events": events})
}

func (h *TerritoryHandler) Decay(c *fiber.Ctx) error {
	res, err := h.territories.Decay(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

type UserHandler struct {
	users *services.UserService
}

func NewUserHandler(users *services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	email, _ := c.Locals("user_email").(string)
	if err := h.users.EnsureProfile(c.Context(), userID, email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	profile, err := h.users.Profile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "profile not found"})
	}
	return c.JSON(profile)
}
