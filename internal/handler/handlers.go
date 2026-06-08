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

// GET /public/projects
func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	featured := c.Query("featured") == "true"
	projects, meta, err := h.svc.List(page, limit, featured, true) // publishedOnly = true
	if err != nil {
		return InternalError(c, "Failed to fetch projects")
	}
	return OKWithMeta(c, projects, meta)
}

// GET /admin/api/projects
func (h *ProjectHandler) AdminListProjects(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	projects, meta, err := h.svc.List(page, limit, false, false)
	if err != nil {
		return InternalError(c, "Failed to fetch projects")
	}
	return OKWithMeta(c, projects, meta)
}

// GET /public/projects/:slug
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	project, err := h.svc.GetBySlug(c.Params("slug"))
	if err != nil {
		return NotFound(c, "Project not found")
	}
	return OK(c, project)
}

// GET /admin/api/projects/:id
func (h *ProjectHandler) AdminGetProject(c *fiber.Ctx) error {
	project, err := h.svc.GetByID(c.Params("id"))
	if err != nil {
		return NotFound(c, "Project not found")
	}
	return OK(c, project)
}

// POST /admin/api/projects
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var req dto.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	project, err := h.svc.Create(&req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, project, "Project created")
}

// PUT /admin/api/projects/:id
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	var req dto.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	project, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Project not found")
	}
	return OK(c, project, "Project updated")
}

// DELETE /admin/api/projects/:id
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Project not found")
	}
	return OK(c, nil, "Project deleted")
}

// ─── Asset Handler (polymorphic: project/experience) ─────────────────────────

type AssetHandler struct{ svc service.AssetService }

func NewAssetHandler(s service.AssetService) *AssetHandler { return &AssetHandler{svc: s} }

// POST /admin/api/assets?owner_type=&owner_id=
func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var req dto.CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	ownerType := c.Query("owner_type") // "project" atau "experience"
	ownerID := c.Query("owner_id")
	if ownerType == "" || ownerID == "" {
		return BadRequest(c, "owner_type and owner_id query params required")
	}
	asset, err := h.svc.Create(ownerType, ownerID, &req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, asset, "Asset created")
}

// GET /admin/api/assets?owner_type=&owner_id=
func (h *AssetHandler) ListAssets(c *fiber.Ctx) error {
	ownerType := c.Query("owner_type")
	ownerID := c.Query("owner_id")
	if ownerType == "" || ownerID == "" {
		return BadRequest(c, "owner_type and owner_id query params required")
	}
	assets, err := h.svc.List(ownerType, ownerID)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return OK(c, assets)
}

// PUT /admin/api/assets/:id
func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	var req dto.UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	asset, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Asset not found")
	}
	return OK(c, asset, "Asset updated")
}

// DELETE /admin/api/assets/:id
func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Asset not found")
	}
	return OK(c, nil, "Asset deleted")
}

// ─── Project Category Handler ────────────────────────────────────────────────

type ProjectCategoryHandler struct{ svc service.ProjectCategoryService }

func NewProjectCategoryHandler(s service.ProjectCategoryService) *ProjectCategoryHandler {
	return &ProjectCategoryHandler{svc: s}
}

// GET /admin/api/project-categories
func (h *ProjectCategoryHandler) ListCategories(c *fiber.Ctx) error {
	cats, err := h.svc.List()
	if err != nil {
		return InternalError(c, "Failed to fetch categories")
	}
	return OK(c, cats)
}

// POST /admin/api/project-categories
func (h *ProjectCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	cat, err := h.svc.Create(&req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, cat, "Category created")
}

// PUT /admin/api/project-categories/:id
func (h *ProjectCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	cat, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Category not found")
	}
	return OK(c, cat, "Category updated")
}

// DELETE /admin/api/project-categories/:id
func (h *ProjectCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Category not found")
	}
	return OK(c, nil, "Category deleted")
}

// ─── Experience Category Handler ─────────────────────────────────────────────

type ExperienceCategoryHandler struct{ svc service.ExperienceCategoryService }

func NewExperienceCategoryHandler(s service.ExperienceCategoryService) *ExperienceCategoryHandler {
	return &ExperienceCategoryHandler{svc: s}
}

// GET /admin/api/experience-categories
func (h *ExperienceCategoryHandler) ListCategories(c *fiber.Ctx) error {
	cats, err := h.svc.List()
	if err != nil {
		return InternalError(c, "Failed to fetch categories")
	}
	return OK(c, cats)
}

// POST /admin/api/experience-categories
func (h *ExperienceCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	cat, err := h.svc.Create(&req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, cat, "Category created")
}

// PUT /admin/api/experience-categories/:id
func (h *ExperienceCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	cat, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Category not found")
	}
	return OK(c, cat, "Category updated")
}

// DELETE /admin/api/experience-categories/:id
func (h *ExperienceCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Category not found")
	}
	return OK(c, nil, "Category deleted")
}

// ─── Stack Category Handler ────────────────────────────────────────────────────

type StackCategoryHandler struct{ svc service.StackCategoryService }

func NewStackCategoryHandler(s service.StackCategoryService) *StackCategoryHandler {
	return &StackCategoryHandler{svc: s}
}

// GET /public/stack-categories?with=items
func (h *StackCategoryHandler) ListCategories(c *fiber.Ctx) error {
	withItems := c.Query("with") == "items"
	categories, err := h.svc.List(withItems)
	if err != nil {
		return InternalError(c, "Failed to fetch stack categories")
	}
	return OK(c, categories)
}

// GET /admin/api/stack-categories
func (h *StackCategoryHandler) AdminListCategories(c *fiber.Ctx) error {
	categories, err := h.svc.List(true)
	if err != nil {
		return InternalError(c, "Failed to fetch stack categories")
	}
	return OK(c, categories)
}

// GET /admin/api/stack-categories/:id
func (h *StackCategoryHandler) GetCategory(c *fiber.Ctx) error {
	cat, err := h.svc.GetByID(c.Params("id"))
	if err != nil {
		return NotFound(c, "Stack category not found")
	}
	return OK(c, cat)
}

// POST /admin/api/stack-categories
func (h *StackCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateStackCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	cat, err := h.svc.Create(&req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, cat, "Stack category created")
}

// PUT /admin/api/stack-categories/:id
func (h *StackCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req dto.UpdateStackCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	cat, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Stack category not found")
	}
	return OK(c, cat, "Stack category updated")
}

// DELETE /admin/api/stack-categories/:id
func (h *StackCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Stack category not found")
	}
	return OK(c, nil, "Stack category deleted")
}

// ─── Stack Item Handler ────────────────────────────────────────────────────────

type StackItemHandler struct{ svc service.StackItemService }

func NewStackItemHandler(s service.StackItemService) *StackItemHandler {
	return &StackItemHandler{svc: s}
}

// GET /public/stack-items?category_id=<uuid>
func (h *StackItemHandler) ListItems(c *fiber.Ctx) error {
	categoryID := c.Query("category_id")
	items, err := h.svc.List(true, categoryID)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return OK(c, items)
}

// GET /admin/api/stack-items?category_id=<uuid>
func (h *StackItemHandler) AdminListItems(c *fiber.Ctx) error {
	categoryID := c.Query("category_id")
	items, err := h.svc.List(false, categoryID)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return OK(c, items)
}

// GET /admin/api/stack-items/:id
func (h *StackItemHandler) GetItem(c *fiber.Ctx) error {
	item, err := h.svc.GetByID(c.Params("id"))
	if err != nil {
		return NotFound(c, "Stack item not found")
	}
	return OK(c, item)
}

// POST /admin/api/stack-items
func (h *StackItemHandler) CreateItem(c *fiber.Ctx) error {
	var req dto.CreateStackItemRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	item, err := h.svc.Create(&req)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, item, "Stack item created")
}

// PUT /admin/api/stack-items/:id
func (h *StackItemHandler) UpdateItem(c *fiber.Ctx) error {
	var req dto.UpdateStackItemRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	item, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Stack item not found")
	}
	return OK(c, item, "Stack item updated")
}

// DELETE /admin/api/stack-items/:id
func (h *StackItemHandler) DeleteItem(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Stack item not found")
	}
	return OK(c, nil, "Stack item deleted")
}

// ─── Skill Handler ─────────────────────────────────────────────────────────────

type SkillHandler struct{ svc service.SkillService }

func NewSkillHandler(s service.SkillService) *SkillHandler { return &SkillHandler{svc: s} }

// GET /public/skills
func (h *SkillHandler) ListSkills(c *fiber.Ctx) error {
	skills, err := h.svc.List(true)
	if err != nil {
		return InternalError(c, "Failed to fetch skills")
	}
	return OK(c, skills)
}

// GET /admin/api/skills
func (h *SkillHandler) AdminListSkills(c *fiber.Ctx) error {
	skills, err := h.svc.List(false)
	if err != nil {
		return InternalError(c, "Failed to fetch skills")
	}
	return OK(c, skills)
}

// GET /public/skills/:id
func (h *SkillHandler) GetSkillDetail(c *fiber.Ctx) error {
	detail, err := h.svc.GetDetail(c.Params("id"))
	if err != nil {
		return NotFound(c, "Skill not found")
	}
	return OK(c, detail)
}

// POST /admin/api/skills
func (h *SkillHandler) CreateSkill(c *fiber.Ctx) error {
	var req dto.CreateSkillRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	skill, err := h.svc.Create(&req)
	if err != nil {
		return InternalError(c, err.Error())
	}
	return Created(c, skill, "Skill created")
}

// PUT /admin/api/skills/:id
func (h *SkillHandler) UpdateSkill(c *fiber.Ctx) error {
	var req dto.UpdateSkillRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	skill, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Skill not found")
	}
	return OK(c, skill, "Skill updated")
}

// DELETE /admin/api/skills/:id
func (h *SkillHandler) DeleteSkill(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Skill not found")
	}
	return OK(c, nil, "Skill deleted")
}

// ─── Experience Handler ────────────────────────────────────────────────────────

type ExperienceHandler struct{ svc service.ExperienceService }

func NewExperienceHandler(s service.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{svc: s}
}

// GET /public/experiences
func (h *ExperienceHandler) ListExperiences(c *fiber.Ctx) error {
	items, err := h.svc.List(true)
	if err != nil {
		return InternalError(c, "Failed to fetch experiences")
	}
	return OK(c, items)
}

// GET /admin/api/experiences
func (h *ExperienceHandler) AdminListExperiences(c *fiber.Ctx) error {
	items, err := h.svc.List(false)
	if err != nil {
		return InternalError(c, "Failed to fetch experiences")
	}
	return OK(c, items)
}

// POST /admin/api/experiences
func (h *ExperienceHandler) CreateExperience(c *fiber.Ctx) error {
	var req dto.CreateExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	item, err := h.svc.Create(&req)
	if err != nil {
		return InternalError(c, err.Error())
	}
	return Created(c, item, "Experience created")
}

// PUT /admin/api/experiences/:id
func (h *ExperienceHandler) UpdateExperience(c *fiber.Ctx) error {
	var req dto.UpdateExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	item, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Experience not found")
	}
	return OK(c, item, "Experience updated")
}

// DELETE /admin/api/experiences/:id
func (h *ExperienceHandler) DeleteExperience(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Experience not found")
	}
	return OK(c, nil, "Experience deleted")
}

// ─── Education Handler ─────────────────────────────────────────────────────────

type EducationHandler struct{ svc service.EducationService }

func NewEducationHandler(s service.EducationService) *EducationHandler {
	return &EducationHandler{svc: s}
}

// GET /public/educations
func (h *EducationHandler) ListEducations(c *fiber.Ctx) error {
	items, err := h.svc.List(true)
	if err != nil {
		return InternalError(c, "Failed to fetch educations")
	}
	return OK(c, items)
}

// GET /admin/api/educations
func (h *EducationHandler) AdminListEducations(c *fiber.Ctx) error {
	items, err := h.svc.List(false)
	if err != nil {
		return InternalError(c, "Failed to fetch educations")
	}
	return OK(c, items)
}

// POST /admin/api/educations
func (h *EducationHandler) CreateEducation(c *fiber.Ctx) error {
	var req dto.CreateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	item, err := h.svc.Create(&req)
	if err != nil {
		return InternalError(c, err.Error())
	}
	return Created(c, item, "Education created")
}

// PUT /admin/api/educations/:id
func (h *EducationHandler) UpdateEducation(c *fiber.Ctx) error {
	var req dto.UpdateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	item, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Education not found")
	}
	return OK(c, item, "Education updated")
}

// DELETE /admin/api/educations/:id
func (h *EducationHandler) DeleteEducation(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Education not found")
	}
	return OK(c, nil, "Education deleted")
}

// ─── Profile Handler ───────────────────────────────────────────────────────────

type ProfileHandler struct{ svc service.ProfileService }

func NewProfileHandler(s service.ProfileService) *ProfileHandler { return &ProfileHandler{svc: s} }

// GET /public/profile
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	profile, err := h.svc.Get()
	if err != nil {
		return NotFound(c, "Profile not found")
	}
	return OK(c, profile)
}

// POST /admin/api/profile
func (h *ProfileHandler) UpsertProfile(c *fiber.Ctx) error {
	var req dto.UpsertProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	profile, err := h.svc.Upsert(&req)
	if err != nil {
		return InternalError(c, err.Error())
	}
	return OK(c, profile, "Profile updated")
}

// ─── Contact Handler ───────────────────────────────────────────────────────────

type ContactHandler struct{ svc service.ContactService }

func NewContactHandler(s service.ContactService) *ContactHandler { return &ContactHandler{svc: s} }

// POST /public/contact
func (h *ContactHandler) SendMessage(c *fiber.Ctx) error {
	var req dto.SendContactRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	if err := h.svc.Save(&req); err != nil {
		return InternalError(c, "Failed to send message")
	}
	return OK(c, nil, "Message sent successfully")
}

// GET /admin/api/contacts
func (h *ContactHandler) AdminListContacts(c *fiber.Ctx) error {
	contacts, err := h.svc.List()
	if err != nil {
		return InternalError(c, "Failed to fetch contacts")
	}
	return OK(c, contacts)
}

// PATCH /admin/api/contacts/:id/read
func (h *ContactHandler) MarkAsRead(c *fiber.Ctx) error {
	if err := h.svc.MarkRead(c.Params("id")); err != nil {
		return NotFound(c, "Contact not found")
	}
	return OK(c, nil, "Marked as read")
}

// ─── Running Activity Handler ────────────────────────────────────────────────

type RunningActivityHandler struct{ svc service.RunningActivityService }

func NewRunningActivityHandler(s service.RunningActivityService) *RunningActivityHandler {
	return &RunningActivityHandler{svc: s}
}

// GET /public/running-activities
func (h *RunningActivityHandler) ListActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	publishedOnly := c.Query("published") != "false" // default true
	activities, meta, err := h.svc.List(publishedOnly, page, limit)
	if err != nil {
		return InternalError(c, "Failed to fetch activities")
	}
	return OKWithMeta(c, activities, meta)
}

// GET /admin/api/running-activities
func (h *RunningActivityHandler) AdminListActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activities, meta, err := h.svc.List(false, page, limit) // admin lihat semua
	if err != nil {
		return InternalError(c, "Failed to fetch activities")
	}
	return OKWithMeta(c, activities, meta)
}

// GET /public/running-activities/:id
func (h *RunningActivityHandler) GetActivity(c *fiber.Ctx) error {
	activity, err := h.svc.GetByID(c.Params("id"))
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
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	activity, err := h.svc.Create(&req)
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
	activity, err := h.svc.Update(c.Params("id"), &req)
	if err != nil {
		return NotFound(c, "Running activity not found")
	}
	return OK(c, activity, "Activity updated")
}

// DELETE /admin/api/running-activities/:id
func (h *RunningActivityHandler) DeleteActivity(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return NotFound(c, "Running activity not found")
	}
	return OK(c, nil, "Activity deleted")
}