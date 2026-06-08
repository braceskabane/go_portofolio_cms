package admin

import (
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupAdmin mounts the built-in admin panel at /admin
func SetupAdmin(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	admin := app.Group("/admin")

	// ── Public admin routes ──
	admin.Get("/login", LoginPage)
	admin.Post("/login", LoginHandler(cfg))
	admin.Get("/logout", LogoutHandler)

	// ── Protected admin routes ──
	protected := admin.Group("", middleware.SessionRequired("/admin/login"))
	protected.Get("/", DashboardPage(db))

	// Projects
	protected.Get("/projects", ProjectsPage(db))
	protected.Get("/projects/new", ProjectFormPage(db, ""))
	protected.Get("/projects/:id/edit", ProjectFormPage(db, ""))
	protected.Post("/projects", CreateProjectHandler(db))
	protected.Post("/projects/:id", UpdateProjectHandler(db))
	protected.Delete("/projects/:id", DeleteProjectHandler(db))

	// Skills
	protected.Get("/skills", SkillsPage(db))
	protected.Get("/skills/new", SkillFormPage(db))
	protected.Get("/skills/:id/edit", SkillFormPage(db))
	protected.Post("/skills", CreateSkillHandler(db))
	protected.Post("/skills/:id", UpdateSkillHandler(db))
	protected.Delete("/skills/:id", DeleteSkillHandler(db))

	// Experiences
	protected.Get("/experiences", ExperiencesPage(db))
	protected.Get("/experiences/new", ExperienceFormPage(db))
	protected.Get("/experiences/:id/edit", ExperienceFormPage(db))
	protected.Post("/experiences", CreateExperienceHandler(db))
	protected.Post("/experiences/:id", UpdateExperienceHandler(db))
	protected.Delete("/experiences/:id", DeleteExperienceHandler(db))

	// Educations
	protected.Get("/educations", EducationsPage(db))
	protected.Get("/educations/new", EducationFormPage(db))
	protected.Get("/educations/:id/edit", EducationFormPage(db))
	protected.Post("/educations", CreateEducationHandler(db))
	protected.Post("/educations/:id", UpdateEducationHandler(db))
	protected.Delete("/educations/:id", DeleteEducationHandler(db))

	// Profile
	protected.Get("/profile", ProfilePage(db))
	protected.Post("/profile", UpsertProfileHandler(db))

	// Contacts
	protected.Get("/contacts", ContactsPage(db))
	protected.Post("/contacts/:id/read", MarkContactReadHandler(db))
	protected.Delete("/contacts/:id", DeleteContactHandler(db))
}
