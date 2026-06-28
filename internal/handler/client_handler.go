package handler

import (
	"strconv"
	"trainers-backend/internal/dto"
	"trainers-backend/internal/models"
	"trainers-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ClientHandler struct {
	clientService      *service.ClientService
	measurementService *service.MeasurementService
	noteService        *service.NoteService
}

func NewClientHandler(cs *service.ClientService, ms *service.MeasurementService, ns *service.NoteService) *ClientHandler {
	return &ClientHandler{clientService: cs, measurementService: ms, noteService: ns}
}

// --- Clients ---
func (h *ClientHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	trainerID := c.Locals("user_id").(string)

	clients, total, err := h.clientService.GetAll(c.Context(), trainerID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}

	return c.JSON(dto.PaginatedResponse{Data: clients, Total: total, Page: page, PageSize: pageSize})
}

func (h *ClientHandler) GetByID(c *fiber.Ctx) error {
	client, err := h.clientService.GetByID(c.Context(), c.Params("id"))
	if err != nil || client == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response{Status: "error", Message: "client not found"})
	}
	return c.JSON(dto.Response{Status: "success", Data: client})
}

func (h *ClientHandler) Create(c *fiber.Ctx) error {
	var req dto.UpdateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}

	client := &models.Client{
		UserID:            c.Params("id"), // В реальном app создаем client для существующего user
		TrainerID:         c.Locals("user_id").(string),
		DateOfBirth:       req.DateOfBirth,
		Gender:            req.Gender,
		HeightCm:          req.HeightCm,
		WeightKg:          req.WeightKg,
		Goal:              req.Goal,
		MedicalConditions: req.MedicalConditions,
	}

	if err := h.clientService.Create(c.Context(), client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.Response{Status: "success", Data: client})
}

func (h *ClientHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}

	client, _ := h.clientService.GetByID(c.Context(), c.Params("id"))
	if client == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response{Status: "error", Message: "not found"})
	}

	client.DateOfBirth = req.DateOfBirth
	client.Gender = req.Gender
	client.HeightCm = req.HeightCm
	client.WeightKg = req.WeightKg
	client.Goal = req.Goal
	client.MedicalConditions = req.MedicalConditions

	if err := h.clientService.Update(c.Context(), client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Data: client})
}

func (h *ClientHandler) Delete(c *fiber.Ctx) error {
	if err := h.clientService.Delete(c.Context(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Message: "deleted"})
}

// --- Measurements ---
func (h *ClientHandler) GetMeasurements(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	data, total, err := h.measurementService.GetByClientID(c.Context(), c.Params("id"), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.PaginatedResponse{Data: data, Total: total, Page: page, PageSize: pageSize})
}

func (h *ClientHandler) CreateMeasurement(c *fiber.Ctx) error {
	var req dto.CreateMeasurementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}

	m := &models.Measurement{
		ClientID: c.Params("id"), MeasurementDate: req.MeasurementDate,
		Weight: req.Weight, BodyFatPercentage: req.BodyFatPercentage,
		MuscleMass: req.MuscleMass, ChestCm: req.ChestCm,
		WaistCm: req.WaistCm, HipsCm: req.HipsCm,
		BicepCm: req.BicepCm, ThighCm: req.ThighCm, Notes: req.Notes,
	}
	if err := h.measurementService.Create(c.Context(), m); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.Response{Status: "success", Data: m})
}

func (h *ClientHandler) DeleteMeasurement(c *fiber.Ctx) error {
	if err := h.measurementService.Delete(c.Context(), c.Params("measurementId")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Message: "deleted"})
}

// --- Notes ---
func (h *ClientHandler) GetNotes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	data, total, err := h.noteService.GetByClientID(c.Context(), c.Params("id"), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.PaginatedResponse{Data: data, Total: total, Page: page, PageSize: pageSize})
}

func (h *ClientHandler) CreateNote(c *fiber.Ctx) error {
	var req dto.CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}

	n := &models.Note{
		ClientID: c.Params("id"), AuthorID: c.Locals("user_id").(string),
		Content: req.Content, IsImportant: req.IsImportant,
	}
	if err := h.noteService.Create(c.Context(), n); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.Response{Status: "success", Data: n})
}

func (h *ClientHandler) UpdateNote(c *fiber.Ctx) error {
	var req dto.UpdateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{Status: "error", Message: "invalid body"})
	}
	n, _ := h.noteService.GetByID(c.Context(), c.Params("noteId")) // Нужен метод GetByID в сервисе
	if n == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response{Status: "error", Message: "not found"})
	}

	n.Content = req.Content
	n.IsImportant = req.IsImportant
	if err := h.noteService.Update(c.Context(), n); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Data: n})
}

func (h *ClientHandler) DeleteNote(c *fiber.Ctx) error {
	if err := h.noteService.Delete(c.Context(), c.Params("noteId")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{Status: "error", Message: err.Error()})
	}
	return c.JSON(dto.Response{Status: "success", Message: "deleted"})
}
