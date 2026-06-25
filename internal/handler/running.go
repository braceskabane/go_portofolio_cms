package handler

import (
	"fmt"
	"io"
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// --------- Running Activity Handler --------------------------------------------------
type RunningActivityHandler struct {
	service       service.RunningActivityService
	geminiService service.GeminiService
}

func NewRunningActivityHandler(s service.RunningActivityService, g service.GeminiService) *RunningActivityHandler {
	return &RunningActivityHandler{service: s, geminiService: g}
}

// GET /admin/api/running-activities
func (h *RunningActivityHandler) ListActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activities, meta, err := h.service.List(true, page, limit)
	if err != nil {
		return InternalError(c, "Failed to fetch activities")
	}
	return OKWithMeta(c, activities, meta)
}

// GET /admin/api/running-activities
func (h *RunningActivityHandler) AdminListActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activities, meta, err := h.service.List(false, page, limit)
	if err != nil {
		return InternalError(c, "Failed to fetch activities")
	}
	return OKWithMeta(c, activities, meta)
}

// GET /admin/api/running-activities/:id
func (h *RunningActivityHandler) GetActivity(c *fiber.Ctx) error {
	activity, err := h.service.GetByID(c.Params("id"))
	if err != nil {
		return NotFound(c, "Running activity not found")
	}
	return OK(c, activity)
}

// POST /admin/api/running-activities
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

// PUT /admin/api/running-activities/:id
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

// POST /admin/api/running-activities/screenshot
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

// ------ Running Analysis Handler --------------------------------------------------

type RunningAnalysisHandler struct {
	analysisService service.RunningAnalysisService
	calendarService service.CalendarService
}

func NewRunningAnalysisHandler(
	analysisService service.RunningAnalysisService,
	calendarService service.CalendarService,
) *RunningAnalysisHandler {
	return &RunningAnalysisHandler{
		analysisService: analysisService,
		calendarService: calendarService,
	}
}

// POST /admin/api/running-analysis/generate
func (h *RunningAnalysisHandler) Generate(c *fiber.Ctx) error {
	var req dto.GenerateRunningAnalysisRequest
	// BodyParser bersifat opsional — request bisa kosong (pakai default)
	_ = c.BodyParser(&req)

	result, err := h.analysisService.GenerateAnalysis(&req)
	if err != nil {
		return InternalError(c, "Gagal generate analisis: "+err.Error())
	}
	return OK(c, result, "Analisis berhasil dibuat")
}

// POST /admin/api/running-analysis/sync-calendar
// Header: X-Calendar-Token: <google_oauth_access_token>
func (h *RunningAnalysisHandler) SyncCalendar(c *fiber.Ctx) error {
	accessToken := c.Get("X-Calendar-Token")
	if accessToken == "" {
		return BadRequest(c, "Header X-Calendar-Token wajib diisi (Google OAuth access token)")
	}

	var req dto.SyncCalendarRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if len(req.Events) == 0 {
		return BadRequest(c, "Minimal 1 event diperlukan")
	}

	result, err := h.calendarService.SyncEvents(accessToken, req.Events)
	if err != nil {
		return InternalError(c, "Gagal sync kalender: "+err.Error())
	}
	return OK(c, result, fmt.Sprintf("%d event berhasil disync ke Google Calendar", result.Synced))
}