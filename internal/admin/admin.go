package admin

import (
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var googleClientID string

func Init(clientID string) {
	googleClientID = clientID
}

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

	// ----- Projects -----
	protected.Get("/projects", ProjectsPage(db))
	protected.Get("/projects/new", ProjectFormPage(db, ""))
	protected.Get("/projects/:id/edit", ProjectFormPage(db, ""))
	protected.Post("/projects", CreateProjectHandler(db))
	protected.Post("/projects/:id", UpdateProjectHandler(db))
	protected.Delete("/projects/:id", DeleteProjectHandler(db))

	// ----- Assets -----
	protected.Get("/assets", AssetsPage(db))
	protected.Get("/assets/new", AssetFormPage(db, ""))
	protected.Get("/assets/:id/edit", AssetFormPage(db, ""))
	protected.Post("/assets", CreateAssetHandler(db))
	protected.Post("/assets/:id", UpdateAssetHandler(db))
	protected.Delete("/assets/:id", DeleteAssetHandler(db))

	// ----- Project Categories -----
	protected.Get("/project-categories", ProjectCategoriesPage(db))
	protected.Get("/project-categories/new", ProjectCategoryFormPage(db, ""))
	protected.Get("/project-categories/:id/edit", ProjectCategoryFormPage(db, ""))
	protected.Post("/project-categories", CreateProjectCategoryHandler(db))
	protected.Post("/project-categories/:id", UpdateProjectCategoryHandler(db))
	protected.Delete("/project-categories/:id", DeleteProjectCategoryHandler(db))

	// ----- Experience Categories -----
	protected.Get("/experience-categories", ExperienceCategoriesPage(db))
	protected.Get("/experience-categories/new", ExperienceCategoryFormPage(db, ""))
	protected.Get("/experience-categories/:id/edit", ExperienceCategoryFormPage(db, ""))
	protected.Post("/experience-categories", CreateExperienceCategoryHandler(db))
	protected.Post("/experience-categories/:id", UpdateExperienceCategoryHandler(db))
	protected.Delete("/experience-categories/:id", DeleteExperienceCategoryHandler(db))

	// ----- Tech Stack Categories -----
	protected.Get("/stack-categories", StackCategoriesPage(db))
	protected.Get("/stack-categories/new", StackCategoryFormPage(db, ""))
	protected.Get("/stack-categories/:id/edit", StackCategoryFormPage(db, ""))
	protected.Post("/stack-categories", CreateStackCategoryHandler(db))
	protected.Post("/stack-categories/:id", UpdateStackCategoryHandler(db))
	protected.Delete("/stack-categories/:id", DeleteStackCategoryHandler(db))

	// ----- Tech Stack Items -----
	protected.Get("/stack-items", StackItemsPage(db))
	protected.Get("/stack-items/new", StackItemFormPage(db, ""))
	protected.Get("/stack-items/:id/edit", StackItemFormPage(db, ""))
	protected.Post("/stack-items", CreateStackItemHandler(db))
	protected.Post("/stack-items/:id", UpdateStackItemHandler(db))
	protected.Delete("/stack-items/:id", DeleteStackItemHandler(db))

	// ----- Skills -----
	protected.Get("/skills", SkillsPage(db))
	protected.Get("/skills/new", SkillFormPage(db))
	protected.Get("/skills/:id/edit", SkillFormPage(db))
	protected.Post("/skills", CreateSkillHandler(db))
	protected.Post("/skills/:id", UpdateSkillHandler(db))
	protected.Delete("/skills/:id", DeleteSkillHandler(db))

	// ----- Experiences -----
	protected.Get("/experiences", ExperiencesPage(db))
	protected.Get("/experiences/new", ExperienceFormPage(db, ""))
	protected.Get("/experiences/:id/edit", ExperienceFormPage(db, ""))
	protected.Post("/experiences", CreateExperienceHandler(db))
	protected.Post("/experiences/:id", UpdateExperienceHandler(db))
	protected.Delete("/experiences/:id", DeleteExperienceHandler(db))

	// ----- Educations -----
	protected.Get("/educations", EducationsPage(db))
	protected.Get("/educations/new", EducationFormPage(db))
	protected.Get("/educations/:id/edit", EducationFormPage(db))
	protected.Post("/educations", CreateEducationHandler(db))
	protected.Post("/educations/:id", UpdateEducationHandler(db))
	protected.Delete("/educations/:id", DeleteEducationHandler(db))

	// ----- Profile -----
	protected.Get("/profile", ProfilePage(db))
	protected.Post("/profile", UpsertProfileHandler(db))

	// ----- Contacts -----
	protected.Get("/contacts", ContactsPage(db))
	protected.Post("/contacts/:id/read", MarkContactReadHandler(db))
	protected.Delete("/contacts/:id", DeleteContactHandler(db))

	// ----- Running Activities -----
	protected.Get("/running-activities", RunningActivitiesPage(db))
	protected.Get("/running-activities/new", RunningActivityFormPage(db, ""))
	protected.Get("/running-activities/:id/edit", RunningActivityFormPage(db, ""))
	protected.Post("/running-activities/preview-screenshots", PreviewScreenshotsHandler(db))
	protected.Post("/running-activities/batch", BatchCreateRunningActivitiesHandler(db)) 
	protected.Post("/running-activities", CreateRunningActivityHandler(db))
	protected.Post("/running-activities/:id", UpdateRunningActivityHandler(db))
	protected.Delete("/running-activities/:id", DeleteRunningActivityHandler(db))

	// ----- Running Analysis -----
	protected.Get("/running-analysis", RunningAnalysisPage(db, cfg))
	protected.Post("/running-analysis/generate", GenerateAnalysisHandler(db, cfg))
}