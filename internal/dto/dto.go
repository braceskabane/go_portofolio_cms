package dto

import "github.com/google/uuid"

// ─── Auth ──────────────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ─── Pagination ────────────────────────────────────────────────────────────────

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ─── Project ───────────────────────────────────────────────────────────────────

type CreateProjectRequest struct {
	Title        string `json:"title" validate:"required,min=1,max=200"`
	Slug         string `json:"slug" validate:"required,min=1,max=250"`
	Description  string `json:"description"`
	Content      string `json:"content"`
	ThumbnailURL string `json:"thumbnail_url" validate:"omitempty,url"`

	// Links
	DemoURL string `json:"demo_url" validate:"omitempty,url"`
	RepoURL string `json:"repo_url" validate:"omitempty,url"`
	DocURL  string `json:"doc_url" validate:"omitempty,url"`

	// Case study
	Problem  string `json:"problem"`
	Solution string `json:"solution"`
	MyRole   string `json:"my_role"`
	Impact   string `json:"impact"`

	StartDate *string `json:"start_date"` // RFC3339
	EndDate   *string `json:"end_date"`   // RFC3339

	IsFeatured  bool `json:"is_featured"`
	IsPublished bool `json:"is_published"`
	SortOrder   int  `json:"sort_order"`

	// Category (optional)
	CategoryID *uuid.UUID `json:"category_id"`

	// Relations — dikirim sebagai array of UUID
	StackItemIDs []uuid.UUID `json:"stack_item_ids"`
	SkillIDs     []uuid.UUID `json:"skill_ids"`
}

type UpdateProjectRequest struct {
	Title        *string `json:"title" validate:"omitempty,min=1,max=200"`
	Slug         *string `json:"slug" validate:"omitempty,min=1,max=250"`
	Description  *string `json:"description"`
	Content      *string `json:"content"`
	ThumbnailURL *string `json:"thumbnail_url" validate:"omitempty,url"`

	DemoURL *string `json:"demo_url" validate:"omitempty,url"`
	RepoURL *string `json:"repo_url" validate:"omitempty,url"`
	DocURL  *string `json:"doc_url" validate:"omitempty,url"`

	Problem  *string `json:"problem"`
	Solution *string `json:"solution"`
	MyRole   *string `json:"my_role"`
	Impact   *string `json:"impact"`

	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`

	IsFeatured  *bool `json:"is_featured"`
	IsPublished *bool `json:"is_published"`
	SortOrder   *int  `json:"sort_order"`

	// Category (opsional, kirim nil untuk tidak ubah, atau pointer ke UUID)
	CategoryID *uuid.UUID `json:"category_id"`

	// nil = jangan update relasi, [] = hapus semua, [ids] = replace
	StackItemIDs *[]uuid.UUID `json:"stack_item_ids"`
	SkillIDs     *[]uuid.UUID `json:"skill_ids"`
}

// ─── Asset (pengganti Media & Document) ────────────────────────────────────────

type CreateAssetRequest struct {
	Type      string `json:"type" validate:"required,oneof=photo video pdf doc"`
	URL       string `json:"url" validate:"required,url"`
	Title     string `json:"title" validate:"omitempty,max=200"`   // untuk dokumen
	Caption   string `json:"caption" validate:"omitempty,max=300"` // untuk foto/video
	SortOrder int    `json:"sort_order"`
}

type UpdateAssetRequest struct {
	Type      *string `json:"type" validate:"omitempty,oneof=photo video pdf doc"`
	URL       *string `json:"url" validate:"omitempty,url"`
	Title     *string `json:"title" validate:"omitempty,max=200"`
	Caption   *string `json:"caption" validate:"omitempty,max=300"`
	SortOrder *int    `json:"sort_order"`
}

// ─── Category (Project & Experience) ───────────────────────────────────────────
// Digunakan oleh ProjectCategoryService dan ExperienceCategoryService

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Slug        string `json:"slug" validate:"required,max=120"`
	SortOrder   int    `json:"sort_order"`
	IsPublished bool   `json:"is_published"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Slug        *string `json:"slug" validate:"omitempty,max=120"`
	SortOrder   *int    `json:"sort_order"`
	IsPublished *bool   `json:"is_published"`
}

// ─── Stack Category ────────────────────────────────────────────────────────────

type CreateStackCategoryRequest struct {
	Name      string `json:"name" validate:"required,max=100"`
	Slug      string `json:"slug" validate:"required,max=120"`
	IconURL   string `json:"icon_url"`
	Color     string `json:"color" validate:"omitempty,max=20"`
	SortOrder int    `json:"sort_order"`
}

type UpdateStackCategoryRequest struct {
	Name      *string `json:"name" validate:"omitempty,max=100"`
	Slug      *string `json:"slug" validate:"omitempty,max=120"`
	IconURL   *string `json:"icon_url"`
	Color     *string `json:"color" validate:"omitempty,max=20"`
	SortOrder *int    `json:"sort_order"`
}

// ─── Stack Item ────────────────────────────────────────────────────────────────

type CreateStackItemRequest struct {
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=100"`
	Slug        string    `json:"slug" validate:"required,max=120"`
	IconURL     string    `json:"icon_url"`
	SortOrder   int       `json:"sort_order"`
	IsPublished bool      `json:"is_published"`
}

type UpdateStackItemRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        *string    `json:"name" validate:"omitempty,max=100"`
	Slug        *string    `json:"slug" validate:"omitempty,max=120"`
	IconURL     *string    `json:"icon_url"`
	SortOrder   *int       `json:"sort_order"`
	IsPublished *bool      `json:"is_published"`
}

// ─── Skill ─────────────────────────────────────────────────────────────────────

type CreateSkillRequest struct {
	Name        string `json:"name"        validate:"required"`
	Category    string `json:"category"`
	IconURL     string `json:"icon_url"`
	Proficiency int    `json:"proficiency"`
	SortOrder   int    `json:"sort_order"`
	IsPublished bool   `json:"is_published"`
}

type UpdateSkillRequest struct {
	Name        *string `json:"name"`
	Category    *string `json:"category"`
	IconURL     *string `json:"icon_url"`
	Proficiency *int    `json:"proficiency"`
	SortOrder   *int    `json:"sort_order"`
	IsPublished *bool   `json:"is_published"`
}

// ─── Experience ────────────────────────────────────────────────────────────────

type CreateExperienceRequest struct {
	Company     string  `json:"company"    validate:"required"`
	Position    string  `json:"position"   validate:"required"`
	Description string  `json:"description"`
	LogoURL     string  `json:"logo_url"`
	CompanyURL  string  `json:"company_url"`
	Location    string  `json:"location"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date"`
	IsCurrent   bool    `json:"is_current"`
	IsPublished bool    `json:"is_published"`
	SortOrder   int     `json:"sort_order"`

	// Category (opsional)
	CategoryID *uuid.UUID `json:"category_id"`

	// Relations
	StackItemIDs []uuid.UUID `json:"stack_item_ids"`
	SkillIDs     []uuid.UUID `json:"skill_ids"`
}

type UpdateExperienceRequest struct {
	Company     *string `json:"company"`
	Position    *string `json:"position"`
	Description *string `json:"description"`
	LogoURL     *string `json:"logo_url"`
	CompanyURL  *string `json:"company_url"`
	Location    *string `json:"location"`
	IsCurrent   *bool   `json:"is_current"`
	IsPublished *bool   `json:"is_published"`
	SortOrder   *int    `json:"sort_order"`

	// Category (opsional, kirim nil untuk tidak ubah, atau pointer ke UUID)
	CategoryID *uuid.UUID `json:"category_id"`

	// Relations (pointer untuk membedakan skip/clear/replace)
	StackItemIDs *[]uuid.UUID `json:"stack_item_ids"`
	SkillIDs     *[]uuid.UUID `json:"skill_ids"`
}

// ─── Education ─────────────────────────────────────────────────────────────────

type CreateEducationRequest struct {
	Institution  string  `json:"institution"    validate:"required"`
	Degree       string  `json:"degree"`
	FieldOfStudy string  `json:"field_of_study"`
	Description  string  `json:"description"`
	LogoURL      string  `json:"logo_url"`
	StartDate    string  `json:"start_date"     validate:"required"`
	EndDate      *string `json:"end_date"`
	GPA          string  `json:"gpa"`
	IsPublished  bool    `json:"is_published"`
	SortOrder    int     `json:"sort_order"`
}

type UpdateEducationRequest struct {
	Institution  *string `json:"institution"`
	Degree       *string `json:"degree"`
	FieldOfStudy *string `json:"field_of_study"`
	IsPublished  *bool   `json:"is_published"`
	SortOrder    *int    `json:"sort_order"`
}

// ─── Profile ───────────────────────────────────────────────────────────────────

type UpsertProfileRequest struct {
	FullName         string `json:"full_name"    validate:"required"`
	Title            string `json:"title"`
	About            string `json:"about"`
	Bio              string `json:"bio"`
	AvatarURL        string `json:"avatar_url"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Location         string `json:"location"`
	GithubURL        string `json:"github_url"`
	LinkedinURL      string `json:"linkedin_url"`
	TwitterURL       string `json:"twitter_url"`
	InstagramURL     string `json:"instagram_url"`
	TiktokURL        string `json:"tiktok_url"`
	StravaURL        string `json:"strava_url"`
	WebsiteURL       string `json:"website_url"`
	ResumeURL        string `json:"resume_url"`
	YearsExperience  int    `json:"years_experience"`
	AvailableForHire bool   `json:"available_for_hire"`
	IsPublished      bool   `json:"is_published"`
}

// ─── Contact ───────────────────────────────────────────────────────────────────

type SendContactRequest struct {
	Name    string `json:"name"    validate:"required,min=2"`
	Email   string `json:"email"   validate:"required,email"`
	Subject string `json:"subject"`
	Message string `json:"message" validate:"required,min=10"`
}

// ─── Running Activity ─────────────────────────────────────────────────────────

type CreateRunningActivityRequest struct {
	Title          string   `json:"title"`
	Notes          string   `json:"notes"`
	MapImageURL    string   `json:"map_image_url"`
	DurationSec    *int     `json:"duration_sec"`      
	DistanceMeters *float64 `json:"distance_meters"`   
	TotalCalories  *float64 `json:"total_calories"`
	ActiveCalories *float64 `json:"active_calories"`   
	AvgPaceSec     *int     `json:"avg_pace_sec"`      
	AvgSpeedKph    *float64 `json:"avg_speed_kph"`     
	AvgCadence     *int     `json:"avg_cadence"`       
	AvgStride      *float64 `json:"avg_stride"`       
	Steps          *int     `json:"steps"`             
	AvgHeartRate   *int     `json:"avg_heart_rate"`    
	IsPublished    bool     `json:"is_published"`
	Date           *string  `json:"date"`
}

type UpdateRunningActivityRequest struct {
	Title          *string  `json:"title"`
	Notes          *string  `json:"notes"`
	MapImageURL    *string  `json:"map_image_url"`
	DurationSec    *int     `json:"duration_sec"`
	DistanceMeters *float64 `json:"distance_meters"`
	TotalCalories  *float64 `json:"total_calories"`
	ActiveCalories *float64 `json:"active_calories"`
	AvgPaceSec     *int     `json:"avg_pace_sec"`
	AvgSpeedKph    *float64 `json:"avg_speed_kph"`
	AvgCadence     *int     `json:"avg_cadence"`
	AvgStride      *float64 `json:"avg_stride"`
	Steps          *int     `json:"steps"`
	AvgHeartRate   *int     `json:"avg_heart_rate"`
	IsPublished    *bool    `json:"is_published"`
	Date           *string  `json:"date"` // RFC3339
}

// ─── Running analysisis sugestion and schedule for google calendar ─────────────────────────────────────────────────────────

// Request body untuk Gemini AI
type GenerateRunningAnalysisRequest struct {
	GoalDescription   string `json:"goal_description"`    // "Ingin bisa lari 10K dalam 2 bulan"
	AvailableDays     []int  `json:"available_days"`      // [1,3,5,6] = Senin,Rabu,Jumat,Sabtu
	PreferredRunTime  string `json:"preferred_run_time"`  // "06:00"
	MaxWeeklyDistKm   float64 `json:"max_weekly_dist_km"` // batas volume (opsional)
}

type SyncCalendarRequest struct {
	Events []CalendarEventDTO `json:"events"`
}

// Gemini Response Contract
type RunningAnalysisResult struct {
	CoachNarrative    string             `json:"coach_narrative"`
	GeneratedAt       string             `json:"generated_at"`
	FitnessAssessment FitnessAssessment  `json:"fitness_assessment"`
	PaceZones         PaceZones          `json:"pace_zones"`
	WeeklyPlan        []DayPlan          `json:"weekly_plan"`
	CalendarEvents    []CalendarEventDTO `json:"calendar_events"`
	Warnings          []RunningWarning   `json:"warnings"`
}

type FitnessAssessment struct {
	Level        string  `json:"level"`         // "beginner" | "intermediate" | "advanced"
	FatigueScore float64     `json:"fatigue_score"` // 0–100
	AerobicBase  string  `json:"aerobic_base"`  // "weak" | "building" | "solid" | "strong"
	Trend        string  `json:"trend"`         // "improving" | "plateau" | "declining"
	CTL          float64 `json:"ctl"`           // chronic training load (42-day avg)
	ATL          float64 `json:"atl"`           // acute training load (7-day avg)
	TSB          float64 `json:"tsb"`           // training stress balance (CTL - ATL)
}

type PaceZones struct {
	Easy      PaceRange `json:"easy"`
	Aerobic   PaceRange `json:"aerobic"`
	Tempo     PaceRange `json:"tempo"`
	Threshold PaceRange `json:"threshold"`
}

type PaceRange struct {
	MinPaceSec int    `json:"min_pace_sec"` // detik/km (konsisten dengan model)
	MaxPaceSec int    `json:"max_pace_sec"`
	Label      string `json:"label"` // "5:30–6:00/km"
}

type DayPlan struct {
	DayOffset   int     `json:"day_offset"`       // 0=hari ini, 1=besok, dst
	DayName     string  `json:"day_name"`         // "Senin"
	SessionType string  `json:"session_type"`     // "easy"|"tempo"|"long"|"rest"|"strength"
	TargetDistKm float64 `json:"target_dist_km"`
	TargetPace  string  `json:"target_pace"`      // "5:45–6:00/km"
	TargetHRZone string `json:"target_hr_zone"`  // "Zone 2 (75–80% HR)"
	Focus       string  `json:"focus"`            // "Fokus cadence 175 spm"
	Rationale   string  `json:"rationale"`        // alasan ilmiah singkat
}

type CalendarEventDTO struct {
	Title       string `json:"title"`
	Date        string `json:"date"`         // "2026-06-28"
	StartTime   string `json:"start_time"`   // "06:00"
	DurationMin int    `json:"duration_min"`
	Description string `json:"description"`  // detail sesi untuk notif
	ColorID     string `json:"color_id"`     // Google Calendar color: "2"=sage, "6"=tangerine, "9"=blueberry
}

type RunningWarning struct {
	Type     string `json:"type"`     // "overtraining"|"injury_risk"|"recovery_needed"
	Message  string `json:"message"`
	Severity string `json:"severity"` // "info"|"warning"|"critical"
}

// ─── Calendar Sync Result ─────────────────────────────────────────────────────

type CalendarSyncResult struct {
	Synced    int                    `json:"synced"`
	Failed    int                    `json:"failed"`
	EventURLs []CalendarEventCreated `json:"event_urls"`
}

type CalendarEventCreated struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	EventID  string `json:"event_id"`
	EventURL string `json:"event_url"`
}
