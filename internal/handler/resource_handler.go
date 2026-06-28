package handler

import (
	"strconv"
	"time"
	"trainers-backend/internal/dto"
	"trainers-backend/internal/models"
	"trainers-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ResourceHandler struct {
	workoutService      *service.WorkoutService
	scheduleService     *service.ScheduleService
	notificationService *service.NotificationService
}

func NewResourceHandler(ws *service.WorkoutService, ss *service.ScheduleService, ns *service.NotificationService) *ResourceHandler {
	return &ResourceHandler{workoutService: ws, scheduleService: ss, notificationService: ns}
}

// --- Workouts ---
func (h *ResourceHandler) GetAllWorkouts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	userID := c.Locals("user_id").(string)
	role := c.Locals("user_role").(string)

	data, total, err := h.workoutService.GetAll(c.Context(), userID, role, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.PaginatedResponse{Data: data, Total: total, Page: page, PageSize: pageSize})
}

func (h *ResourceHandler) GetWorkoutByID(c *fiber.Ctx) error {
	w, err := h.workoutService.GetByID(c.Context(), c.Params("id"))
	if err != nil || w == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response{Status: "error", Message: "not found"})
	}
	return c.JSON(dto.Response{Status: "success", Data: w})
}

func (h *ResourceHandler) CreateWorkout(c *fiber.Ctx) error {
	var req dto.CreateWorkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}
	schedAt, _ := time.Parse(time.RFC3339, req.ScheduledAt)

	w := &models.Workout{
		ClientID: req.ClientID, TrainerID: c.Locals("user_id").(string),
		Title: req.Title, Description: req.Description, WorkoutType: req.WorkoutType,
		ScheduledAt: schedAt, DurationMinutes: req.DurationMinutes,
		Exercises: req.Exercises, Notes: req.Notes,
	}
	if err := h.workoutService.Create(c.Context(), w); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.Response{Status: "success", Data: w})
}

func (h *ResourceHandler) UpdateWorkout(c *fiber.Ctx) error {
	var req dto.UpdateWorkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}
	w, _ := h.workoutService.GetByID(c.Context(), c.Params("id"))
	if w == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response{Status: "error", Message: "not found"})
	}

	w.Title = req.Title
	w.Description = req.Description
	w.WorkoutType = req.WorkoutType
	w.DurationMinutes = req.DurationMinutes
	w.Status = req.Status
	w.Exercises = req.Exercises
	w.Notes = req.Notes
	if req.ScheduledAt != "" {
		w.ScheduledAt, _ = time.Parse(time.RFC3339, req.ScheduledAt)
	}

	if err := h.workoutService.Update(c.Context(), w); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Data: w})
}

func (h *ResourceHandler) DeleteWorkout(c *fiber.Ctx) error {
	if err := h.workoutService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Message: "deleted"})
}

// --- Schedule ---
func (h *ResourceHandler) GetSchedule(c *fiber.Ctx) error {
	trainerID := c.Locals("user_id").(string)
	data, err := h.scheduleService.GetAll(c.Context(), trainerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Data: data})
}

func (h *ResourceHandler) CreateSchedule(c *fiber.Ctx) error {
	var req dto.CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}
	s := &models.Schedule{
		TrainerID: c.Locals("user_id").(string), DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime, EndTime: req.EndTime, IsAvailable: req.IsAvailable,
	}
	if err := h.scheduleService.Create(c.Context(), s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.Response{Status: "success", Data: s})
}

func (h *ResourceHandler) UpdateSchedule(c *fiber.Ctx) error {
	var req dto.UpdateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}
	s, _ := h.scheduleService.GetByID(c.Context(), c.Params("id"))
	if s == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response{Status: "error", Message: "not found"})
	}

	s.DayOfWeek = req.DayOfWeek
	s.StartTime = req.StartTime
	s.EndTime = req.EndTime
	s.IsAvailable = req.IsAvailable
	if err := h.scheduleService.Update(c.Context(), s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Data: s})
}

func (h *ResourceHandler) DeleteSchedule(c *fiber.Ctx) error {
	if err := h.scheduleService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Message: "deleted"})
}

// --- Notifications ---
func (h *ResourceHandler) GetNotifications(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	userID := c.Locals("user_id").(string)

	data, total, err := h.notificationService.GetByUserID(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.PaginatedResponse{Data: data, Total: total, Page: page, PageSize: pageSize})
}

func (h *ResourceHandler) MarkNotificationRead(c *fiber.Ctx) error {
	if err := h.notificationService.MarkAsRead(c.Context(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Message: "marked as read"})
}
