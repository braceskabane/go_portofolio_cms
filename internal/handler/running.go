package handler

import (
	"io"
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RunningActivityHandler struct {
	service       service.RunningActivityService
	geminiService service.GeminiService
}

func NewRunningActivityHandler(s service.RunningActivityService, g service.GeminiService) *RunningActivityHandler {
	return &RunningActivityHandler{service: s, geminiService: g}
}

func (h *RunningActivityHandler) ListActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activities, meta, err := h.service.List(true, page, limit)
	if err != nil {
		return InternalError(c, "Failed to fetch activities")
	}
	return OKWithMeta(c, activities, meta)
}

func (h *RunningActivityHandler) AdminListActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activities, meta, err := h.service.List(false, page, limit)
	if err != nil {
		return InternalError(c, "Failed to fetch activities")
	}
	return OKWithMeta(c, activities, meta)
}

func (h *RunningActivityHandler) GetActivity(c *fiber.Ctx) error {
	activity, err := h.service.GetByID(c.Params("id"))
	if err != nil {
		return NotFound(c, "Running activity not found")
	}
	return OK(c, activity)
}

func (h *RunningActivityHandler) CreateActivity(c *fiber.Ctx) error {
	var req dto.CreateRunningActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	activity, err := h.service.Create(&req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, activity, "Activity created")
}

func (h *RunningActivityHandler) UpdateActivity(c *fiber.Ctx) error {
	var req dto.UpdateRunningActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	activity, err := h.service.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Running activity not found")
	}
	return OK(c, activity, "Activity updated")
}

func (h *RunningActivityHandler) DeleteActivity(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Running activity not found")
	}
	return OK(c, nil, "Activity deleted")
}

func (h *RunningActivityHandler) CreateActivityFromScreenshot(c *fiber.Ctx) error {
	file, err := c.FormFile("screenshot")
	if err != nil {
		return BadRequest(c, "Field 'screenshot' wajib diisi")
	}

	mimeType := file.Header.Get("Content-Type")
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowed[mimeType] {
		return BadRequest(c, "Hanya file JPEG, PNG, atau WebP yang diizinkan")
	}

	if file.Size > 10*1024*1024 {
		return BadRequest(c, "Ukuran file maksimal 10MB")
	}

	f, err := file.Open()
	if err != nil {
		return InternalError(c, "Gagal membaca file")
	}
	defer f.Close()

	imageBytes, err := io.ReadAll(f)
	if err != nil {
		return InternalError(c, "Gagal membaca file")
	}

	extracted, err := h.geminiService.ExtractRunningActivity(imageBytes, mimeType)
	if err != nil {
		return InternalError(c, "Gagal mengekstrak data: "+err.Error())
	}

	activity, err := h.service.Create(extracted)
	if err != nil {
		return InternalError(c, "Gagal menyimpan aktivitas")
	}

	return Created(c, activity, "Aktivitas berhasil dibuat dari screenshot")
}