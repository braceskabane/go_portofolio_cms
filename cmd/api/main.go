package main

import (
	"log"
	_ "portfolio-cms/docs"

	"portfolio-cms/internal/admin"
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/database"
	"portfolio-cms/internal/handler"
	"portfolio-cms/internal/middleware"
	"portfolio-cms/internal/service"

	"github.com/gofiber/fiber/v2"
)

// @title           Portfolio CMS API
// @version         1.0
// @description     REST API for Portfolio CMS — Go Fiber + GORM + PostgreSQL
// @contact.name    API Support
// @license.name    MIT
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// ── 1. Load config ────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── 2. Connect database & auto migrate ────────────────────────────────────
	db := database.Connect(cfg) // AutoMigrate sudah dijalankan di Connect

	// ── 3. Init Fiber app ─────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// ── 4. Global middleware ──────────────────────────────────────────────────
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// ── 5. Mount custom admin panel at /admin ────────────────────────────────
	admin.Init(cfg.Google.ClientID)
	admin.SetupAdmin(app, cfg, db)

	// ── 6. Wire services ──────────────────────────────────────────────────────
	authSvc                := service.NewAuthService(db)
	projectSvc             := service.NewProjectService(db)
	assetSvc               := service.NewAssetService(db)
	projectCategorySvc     := service.NewProjectCategoryService(db)
	experienceCategorySvc  := service.NewExperienceCategoryService(db)
	stackCategorySvc       := service.NewStackCategoryService(db)
	stackItemSvc           := service.NewStackItemService(db)
	skillSvc               := service.NewSkillService(db)
	experienceSvc          := service.NewExperienceService(db)
	educationSvc           := service.NewEducationService(db)
	profileSvc             := service.NewProfileService(db)
	contactSvc             := service.NewContactService(db)
	runningActivitySvc     := service.NewRunningActivityService(db)
	geminiSvc             := service.NewGeminiService(cfg.Gemini.APIKey, cfg.Gemini.Model)
	runningAnalysisSvc := service.NewRunningAnalysisService(db, cfg.Gemini.APIKey, cfg.Gemini.Model)
	calendarSvc        := service.NewCalendarService()

	// ── 7. Wire handlers ──────────────────────────────────────────────────────
	authHandler                := handler.NewAuthHandler(authSvc)
	projectHandler             := handler.NewProjectHandler(projectSvc)
	assetHandler               := handler.NewAssetHandler(assetSvc)
	projectCategoryHandler     := handler.NewProjectCategoryHandler(projectCategorySvc)
	experienceCategoryHandler  := handler.NewExperienceCategoryHandler(experienceCategorySvc)
	stackCategoryHandler       := handler.NewStackCategoryHandler(stackCategorySvc)
	stackItemHandler           := handler.NewStackItemHandler(stackItemSvc)
	skillHandler               := handler.NewSkillHandler(skillSvc)
	experienceHandler          := handler.NewExperienceHandler(experienceSvc)
	educationHandler           := handler.NewEducationHandler(educationSvc)
	profileHandler             := handler.NewProfileHandler(profileSvc)
	contactHandler             := handler.NewContactHandler(contactSvc)
	runningActivityHandler := handler.NewRunningActivityHandler(runningActivitySvc, geminiSvc) 
	runningAnalysisHandler := handler.NewRunningAnalysisHandler(runningAnalysisSvc, calendarSvc)

	// ── 8. Register routes ────────────────────────────────────────────────────
	router := NewRouter(
		app,
		authHandler,
		projectHandler,
		assetHandler,
		projectCategoryHandler,
		experienceCategoryHandler,
		stackCategoryHandler,
		stackItemHandler,
		skillHandler,
		experienceHandler,
		educationHandler,
		profileHandler,
		contactHandler,
		runningActivityHandler,
		runningAnalysisHandler,
	)
	router.Setup()

	// ── 9. Start server ───────────────────────────────────────────────────────
	log.Printf("🚀 Server running on port %s", cfg.App.Port)
	log.Printf("📋 Admin panel: http://localhost:%s/admin", cfg.App.Port)
	log.Printf("📚 API docs:    http://localhost:%s/docs/index.html", cfg.App.Port)

	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// customErrorHandler returns consistent JSON error responses
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}