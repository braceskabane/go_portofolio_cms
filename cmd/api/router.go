package main

import (
	"portfolio-cms/internal/handler"
	"portfolio-cms/internal/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

type Router struct {
	app *fiber.App

	// Handlers
	auth                 *handler.AuthHandler
	project              *handler.ProjectHandler
	asset                *handler.AssetHandler
	projectCategory      *handler.ProjectCategoryHandler
	experienceCategory   *handler.ExperienceCategoryHandler
	stackCategory        *handler.StackCategoryHandler
	stackItem            *handler.StackItemHandler
	skill                *handler.SkillHandler
	experience           *handler.ExperienceHandler
	education            *handler.EducationHandler
	profile              *handler.ProfileHandler
	contact              *handler.ContactHandler
	runningActivity      *handler.RunningActivityHandler
	runningAnalysis      *handler.RunningAnalysisHandler
}

func NewRouter(
	app *fiber.App,
	auth *handler.AuthHandler,
	project *handler.ProjectHandler,
	asset *handler.AssetHandler,
	projectCategory *handler.ProjectCategoryHandler,
	experienceCategory *handler.ExperienceCategoryHandler,
	stackCategory *handler.StackCategoryHandler,
	stackItem *handler.StackItemHandler,
	skill *handler.SkillHandler,
	experience *handler.ExperienceHandler,
	education *handler.EducationHandler,
	profile *handler.ProfileHandler,
	contact *handler.ContactHandler,
	runningActivity *handler.RunningActivityHandler,
	runningAnalysis *handler.RunningAnalysisHandler,
) *Router {
	return &Router{
		app:                 app,
		auth:                auth,
		project:             project,
		asset:               asset,
		projectCategory:     projectCategory,
		experienceCategory:  experienceCategory,
		stackCategory:       stackCategory,
		stackItem:           stackItem,
		skill:               skill,
		experience:          experience,
		education:           education,
		profile:             profile,
		contact:             contact,
		runningActivity:     runningActivity,
		runningAnalysis:     runningAnalysis,
	}
}

func (r *Router) Setup() {
	// Swagger docs
	r.app.Get("/docs/*", fiberSwagger.HandlerDefault)

	// Health check
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "portfolio-cms"})
	})

	// ── API v1 ────────────────────────────────────────────────────────────────
	v1 := r.app.Group("/api/v1")

	// Auth
	auth := v1.Group("/auth")
	auth.Post("/register", r.auth.Register)
	auth.Post("/login", r.auth.Login)
	auth.Post("/refresh", r.auth.RefreshToken)
	auth.Get("/me", middleware.JWTProtected(), r.auth.Me)

	// ── Public routes (for Nuxt frontend) ────────────────────────────────────
	public := v1.Group("/public")

	public.Get("/profile", r.profile.GetProfile)
	public.Get("/projects", r.project.ListProjects)
	public.Get("/projects/:slug", r.project.GetProject)
	public.Get("/skills", r.skill.ListSkills)
	public.Get("/skills/:id", r.skill.GetSkillDetail)
	public.Get("/experiences", r.experience.ListExperiences)
	public.Get("/educations", r.education.ListEducations)
	public.Post("/contact", r.contact.SendMessage)

	// Public: Tech Stack
	public.Get("/stack-categories", r.stackCategory.ListCategories)
	public.Get("/stack-items", r.stackItem.ListItems)

	// Public: Running Activities
	public.Get("/running-activities", r.runningActivity.ListActivities)
	public.Get("/running-activities/:id", r.runningActivity.GetActivity)

	// ── Protected API routes (for admin dashboard) ───────────────────────────
	protected := v1.Group("/admin/api", middleware.JWTProtected(), middleware.RoleRequired("admin", "superadmin"))

	// Projects
	protected.Get("/projects", r.project.AdminListProjects)
	protected.Post("/projects", r.project.CreateProject)
	protected.Get("/projects/:id", r.project.AdminGetProject)
	protected.Put("/projects/:id", r.project.UpdateProject)
	protected.Delete("/projects/:id", r.project.DeleteProject)

	// Assets (polymorphic)
	protected.Post("/assets", r.asset.CreateAsset)
	protected.Get("/assets", r.asset.ListAssets)
	protected.Put("/assets/:id", r.asset.UpdateAsset)
	protected.Delete("/assets/:id", r.asset.DeleteAsset)

	// Project Categories
	protected.Get("/project-categories", r.projectCategory.ListCategories)
	protected.Post("/project-categories", r.projectCategory.CreateCategory)
	protected.Put("/project-categories/:id", r.projectCategory.UpdateCategory)
	protected.Delete("/project-categories/:id", r.projectCategory.DeleteCategory)

	// Experience Categories
	protected.Get("/experience-categories", r.experienceCategory.ListCategories)
	protected.Post("/experience-categories", r.experienceCategory.CreateCategory)
	protected.Put("/experience-categories/:id", r.experienceCategory.UpdateCategory)
	protected.Delete("/experience-categories/:id", r.experienceCategory.DeleteCategory)

	// Tech Stack Categories
	protected.Get("/stack-categories", r.stackCategory.AdminListCategories)
	protected.Get("/stack-categories/:id", r.stackCategory.GetCategory)
	protected.Post("/stack-categories", r.stackCategory.CreateCategory)
	protected.Put("/stack-categories/:id", r.stackCategory.UpdateCategory)
	protected.Delete("/stack-categories/:id", r.stackCategory.DeleteCategory)

	// Tech Stack Items
	protected.Get("/stack-items", r.stackItem.AdminListItems)
	protected.Get("/stack-items/:id", r.stackItem.GetItem)
	protected.Post("/stack-items", r.stackItem.CreateItem)
	protected.Put("/stack-items/:id", r.stackItem.UpdateItem)
	protected.Delete("/stack-items/:id", r.stackItem.DeleteItem)

	// Skills
	protected.Get("/skills", r.skill.AdminListSkills)
	protected.Post("/skills", r.skill.CreateSkill)
	protected.Put("/skills/:id", r.skill.UpdateSkill)
	protected.Delete("/skills/:id", r.skill.DeleteSkill)

	// Experiences
	protected.Get("/experiences", r.experience.AdminListExperiences)
	protected.Post("/experiences", r.experience.CreateExperience)
	protected.Put("/experiences/:id", r.experience.UpdateExperience)
	protected.Delete("/experiences/:id", r.experience.DeleteExperience)

	// Educations
	protected.Get("/educations", r.education.AdminListEducations)
	protected.Post("/educations", r.education.CreateEducation)
	protected.Put("/educations/:id", r.education.UpdateEducation)
	protected.Delete("/educations/:id", r.education.DeleteEducation)

	// Profile
	protected.Get("/profile", r.profile.GetProfile)
	protected.Post("/profile", r.profile.UpsertProfile)

	// Contacts
	protected.Get("/contacts", r.contact.AdminListContacts)
	protected.Patch("/contacts/:id/read", r.contact.MarkAsRead)

	// Running Activities
	protected.Get("/running-activities", r.runningActivity.AdminListActivities)
	protected.Post("/running-activities/screenshot", r.runningActivity.CreateActivityFromScreenshot) 
	protected.Post("/running-activities", r.runningActivity.CreateActivity)
	protected.Put("/running-activities/:id", r.runningActivity.UpdateActivity)
	protected.Delete("/running-activities/:id", r.runningActivity.DeleteActivity)

	// Running Analysis
	protected.Post("/running-analysis/generate",        r.runningAnalysis.Generate)
	protected.Post("/running-analysis/sync-calendar",   r.runningAnalysis.SyncCalendar)
}