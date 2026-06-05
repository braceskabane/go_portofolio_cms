package handler

import (
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ─── Project Handler ───────────────────────────────────────────────────────────

type ProjectHandler struct{ svc service.ProjectService }
func NewProjectHandler(s service.ProjectService) *ProjectHandler { return &ProjectHandler{svc: s} }

func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	projects, meta, err := h.svc.List(page, limit, false, true)
	if err != nil { return InternalError(c, "Failed to fetch projects") }
	return OKWithMeta(c, projects, meta)
}

func (h *ProjectHandler) AdminListProjects(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	projects, meta, err := h.svc.List(page, limit, false, false)
	if err != nil { return InternalError(c, "Failed to fetch projects") }
	return OKWithMeta(c, projects, meta)
}

func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	project, err := h.svc.GetBySlug(c.Params("slug"))
	if err != nil { return NotFound(c, "Project not found") }
	return OK(c, project)
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var req dto.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	if err := validate.Struct(req); err != nil { return BadRequest(c, err.Error()) }
	project, err := h.svc.Create(&req)
	if err != nil { return InternalError(c, err.Error()) }
	return Created(c, project, "Project created")
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	var req dto.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	project, err := h.svc.Update(c.Params("id"), &req)
	if err != nil { return NotFound(c, "Project not found") }
	return OK(c, project, "Project updated")
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil { return NotFound(c, "Project not found") }
	return OK(c, nil, "Project deleted")
}

// ─── Skill Handler ─────────────────────────────────────────────────────────────

type SkillHandler struct{ svc service.SkillService }
func NewSkillHandler(s service.SkillService) *SkillHandler { return &SkillHandler{svc: s} }

func (h *SkillHandler) ListSkills(c *fiber.Ctx) error {
	skills, err := h.svc.List(true)
	if err != nil { return InternalError(c, "Failed to fetch skills") }
	return OK(c, skills)
}

func (h *SkillHandler) AdminListSkills(c *fiber.Ctx) error {
	skills, err := h.svc.List(false)
	if err != nil { return InternalError(c, "Failed to fetch skills") }
	return OK(c, skills)
}

func (h *SkillHandler) CreateSkill(c *fiber.Ctx) error {
	var req dto.CreateSkillRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	if err := validate.Struct(req); err != nil { return BadRequest(c, err.Error()) }
	skill, err := h.svc.Create(&req)
	if err != nil { return InternalError(c, err.Error()) }
	return Created(c, skill, "Skill created")
}

func (h *SkillHandler) UpdateSkill(c *fiber.Ctx) error {
	var req dto.UpdateSkillRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	skill, err := h.svc.Update(c.Params("id"), &req)
	if err != nil { return NotFound(c, "Skill not found") }
	return OK(c, skill, "Skill updated")
}

func (h *SkillHandler) DeleteSkill(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil { return NotFound(c, "Skill not found") }
	return OK(c, nil, "Skill deleted")
}

// ─── Experience Handler ────────────────────────────────────────────────────────

type ExperienceHandler struct{ svc service.ExperienceService }
func NewExperienceHandler(s service.ExperienceService) *ExperienceHandler { return &ExperienceHandler{svc: s} }

func (h *ExperienceHandler) ListExperiences(c *fiber.Ctx) error {
	items, err := h.svc.List(true)
	if err != nil { return InternalError(c, "Failed to fetch experiences") }
	return OK(c, items)
}

func (h *ExperienceHandler) AdminListExperiences(c *fiber.Ctx) error {
	items, err := h.svc.List(false)
	if err != nil { return InternalError(c, "Failed to fetch experiences") }
	return OK(c, items)
}

func (h *ExperienceHandler) CreateExperience(c *fiber.Ctx) error {
	var req dto.CreateExperienceRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	if err := validate.Struct(req); err != nil { return BadRequest(c, err.Error()) }
	item, err := h.svc.Create(&req)
	if err != nil { return InternalError(c, err.Error()) }
	return Created(c, item, "Experience created")
}

func (h *ExperienceHandler) UpdateExperience(c *fiber.Ctx) error {
	var req dto.UpdateExperienceRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	item, err := h.svc.Update(c.Params("id"), &req)
	if err != nil { return NotFound(c, "Experience not found") }
	return OK(c, item, "Experience updated")
}

func (h *ExperienceHandler) DeleteExperience(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil { return NotFound(c, "Experience not found") }
	return OK(c, nil, "Experience deleted")
}

// ─── Education Handler ─────────────────────────────────────────────────────────

type EducationHandler struct{ svc service.EducationService }
func NewEducationHandler(s service.EducationService) *EducationHandler { return &EducationHandler{svc: s} }

func (h *EducationHandler) ListEducations(c *fiber.Ctx) error {
	items, err := h.svc.List(true)
	if err != nil { return InternalError(c, "Failed to fetch educations") }
	return OK(c, items)
}

func (h *EducationHandler) CreateEducation(c *fiber.Ctx) error {
	var req dto.CreateEducationRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	if err := validate.Struct(req); err != nil { return BadRequest(c, err.Error()) }
	item, err := h.svc.Create(&req)
	if err != nil { return InternalError(c, err.Error()) }
	return Created(c, item, "Education created")
}

func (h *EducationHandler) UpdateEducation(c *fiber.Ctx) error {
	var req dto.UpdateEducationRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	item, err := h.svc.Update(c.Params("id"), &req)
	if err != nil { return NotFound(c, "Education not found") }
	return OK(c, item, "Education updated")
}

func (h *EducationHandler) DeleteEducation(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil { return NotFound(c, "Education not found") }
	return OK(c, nil, "Education deleted")
}

// ─── Profile Handler ───────────────────────────────────────────────────────────

type ProfileHandler struct{ svc service.ProfileService }
func NewProfileHandler(s service.ProfileService) *ProfileHandler { return &ProfileHandler{svc: s} }

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	profile, err := h.svc.Get()
	if err != nil { return NotFound(c, "Profile not found") }
	return OK(c, profile)
}

func (h *ProfileHandler) UpsertProfile(c *fiber.Ctx) error {
	var req dto.UpsertProfileRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	if err := validate.Struct(req); err != nil { return BadRequest(c, err.Error()) }
	profile, err := h.svc.Upsert(&req)
	if err != nil { return InternalError(c, err.Error()) }
	return OK(c, profile, "Profile updated")
}

// ─── Contact Handler ───────────────────────────────────────────────────────────

type ContactHandler struct{ svc service.ContactService }
func NewContactHandler(s service.ContactService) *ContactHandler { return &ContactHandler{svc: s} }

func (h *ContactHandler) SendMessage(c *fiber.Ctx) error {
	var req dto.SendContactRequest
	if err := c.BodyParser(&req); err != nil { return BadRequest(c, "Invalid request body") }
	if err := validate.Struct(req); err != nil { return BadRequest(c, err.Error()) }
	if err := h.svc.Save(&req); err != nil { return InternalError(c, "Failed to send message") }
	return OK(c, nil, "Message sent successfully")
}

func (h *ContactHandler) AdminListContacts(c *fiber.Ctx) error {
	contacts, err := h.svc.List()
	if err != nil { return InternalError(c, "Failed to fetch contacts") }
	return OK(c, contacts)
}

func (h *ContactHandler) MarkAsRead(c *fiber.Ctx) error {
	if err := h.svc.MarkRead(c.Params("id")); err != nil { return NotFound(c, "Contact not found") }
	return OK(c, nil, "Marked as read")
}
