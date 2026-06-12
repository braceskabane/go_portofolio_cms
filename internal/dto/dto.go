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