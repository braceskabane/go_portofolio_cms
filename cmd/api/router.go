package main

import (
	"portfolio-cms/internal/handler"
	"portfolio-cms/internal/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

type Router struct {
	app        *fiber.App
	auth       *handler.AuthHandler
	project    *handler.ProjectHandler
	skill      *handler.SkillHandler
	experience *handler.ExperienceHandler
	education  *handler.EducationHandler
	profile    *handler.ProfileHandler
	contact    *handler.ContactHandler
}

func NewRouter(
	app *fiber.App,
	auth *handler.AuthHandler,
	project *handler.ProjectHandler,
	skill *handler.SkillHandler,
	experience *handler.ExperienceHandler,
	education *handler.EducationHandler,
	profile *handler.ProfileHandler,
	contact *handler.ContactHandler,
) *Router {
	return &Router{
		app:        app,
		auth:       auth,
		project:    project,
		skill:      skill,
		experience: experience,
		education:  education,
		profile:    profile,
		contact:    contact,
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
	public.Get("/experiences", r.experience.ListExperiences)
	public.Get("/educations", r.education.ListEducations)
	public.Post("/contact", r.contact.SendMessage)

	// ── Protected API routes (for custom admin or testing) ───────────────────
	protected := v1.Group("/admin/api", middleware.JWTProtected(), middleware.RoleRequired("admin", "superadmin"))

	// Projects
	protected.Get("/projects", r.project.AdminListProjects)
	protected.Post("/projects", r.project.CreateProject)
	protected.Put("/projects/:id", r.project.UpdateProject)
	protected.Delete("/projects/:id", r.project.DeleteProject)

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
	protected.Get("/educations", r.education.ListEducations)
	protected.Post("/educations", r.education.CreateEducation)
	protected.Put("/educations/:id", r.education.UpdateEducation)
	protected.Delete("/educations/:id", r.education.DeleteEducation)

	// Profile
	protected.Get("/profile", r.profile.GetProfile)
	protected.Post("/profile", r.profile.UpsertProfile)

	// Contacts
	protected.Get("/contacts", r.contact.AdminListContacts)
	protected.Patch("/contacts/:id/read", r.contact.MarkAsRead)
}
