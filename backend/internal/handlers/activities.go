package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/territory-run/api/internal/middleware"
	"github.com/territory-run/api/internal/models"
	"github.com/territory-run/api/internal/services"
)

type ActivityHandler struct {
	activities *services.ActivityService
	users      *services.UserService
}

func NewActivityHandler(activities *services.ActivityService, users *services.UserService) *ActivityHandler {
	return &ActivityHandler{activities: activities, users: users}
}

func (h *ActivityHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	email, _ := c.Locals("user_email").(string)
	if err := h.users.EnsureProfile(c.Context(), userID, email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var req models.CreateActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.IdempotencyKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "idempotency_key required"})
	}

	id, err := h.activities.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"activity_id": id})
}

func (h *ActivityHandler) UploadPoints(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	activityID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid activity id"})
	}

	var req models.UploadPointsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if len(req.Points) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "points required"})
	}

	if err := h.activities.AppendPoints(c.Context(), activityID, userID, req.Points); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "count": len(req.Points)})
}

func (h *ActivityHandler) Complete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	activityID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid activity id"})
	}

	var req models.CompleteActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	resp, err := h.activities.Complete(c.Context(), activityID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

func (h *ActivityHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	items, err := h.users.ListActivities(c.Context(), userID, 20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"activities": items})
}
