package dto

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
	Title        string `json:"title"         validate:"required,min=2,max=200"`
	Slug         string `json:"slug"          validate:"required"`
	Description  string `json:"description"`
	Content      string `json:"content"`
	ThumbnailURL string `json:"thumbnail_url"`
	DemoURL      string `json:"demo_url"`
	RepoURL      string `json:"repo_url"`
	TechStack    string `json:"tech_stack"`
	IsFeatured   bool   `json:"is_featured"`
	IsPublished  bool   `json:"is_published"`
	SortOrder    int    `json:"sort_order"`
}

type UpdateProjectRequest struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	Content      *string `json:"content"`
	ThumbnailURL *string `json:"thumbnail_url"`
	DemoURL      *string `json:"demo_url"`
	RepoURL      *string `json:"repo_url"`
	TechStack    *string `json:"tech_stack"`
	IsFeatured   *bool   `json:"is_featured"`
	IsPublished  *bool   `json:"is_published"`
	SortOrder    *int    `json:"sort_order"`
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
}

type UpdateExperienceRequest struct {
	Company     *string `json:"company"`
	Position    *string `json:"position"`
	Description *string `json:"description"`
	LogoURL     *string `json:"logo_url"`
	IsCurrent   *bool   `json:"is_current"`
	IsPublished *bool   `json:"is_published"`
	SortOrder   *int    `json:"sort_order"`
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
	FullName    string `json:"full_name"    validate:"required"`
	Title       string `json:"title"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Location    string `json:"location"`
	GithubURL   string `json:"github_url"`
	LinkedinURL string `json:"linkedin_url"`
	TwitterURL  string `json:"twitter_url"`
	WebsiteURL  string `json:"website_url"`
	ResumeURL   string `json:"resume_url"`
	IsPublished bool   `json:"is_published"`
}

// ─── Contact ───────────────────────────────────────────────────────────────────

type SendContactRequest struct {
	Name    string `json:"name"    validate:"required,min=2"`
	Email   string `json:"email"   validate:"required,email"`
	Subject string `json:"subject"`
	Message string `json:"message" validate:"required,min=10"`
}
