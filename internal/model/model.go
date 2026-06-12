package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Base Model ────────────────────────────────────────────────────────────────

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// ─── User (Admin) ──────────────────────────────────────────────────────────────

type User struct {
	Base
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"size:20;default:'admin'" json:"role"` // admin | superadmin
}

// ─── Profile (single row - owner info) ────────────────────────────────────────

type Profile struct {
	Base
	FullName         string `gorm:"size:150;not null" json:"full_name"`
	Title            string `gorm:"size:200" json:"title"`
	About            string `gorm:"type:text" json:"about"`
	Bio              string `gorm:"type:text" json:"bio"`
	AvatarURL        string `gorm:"size:500" json:"avatar_url"`
	Location         string `gorm:"size:150" json:"location"`
	Email            string `gorm:"size:150" json:"email"`
	Phone            string `gorm:"size:30" json:"phone"`
	GithubURL        string `gorm:"size:300" json:"github_url"`
	LinkedinURL      string `gorm:"size:300" json:"linkedin_url"`
	TwitterURL       string `gorm:"size:300" json:"twitter_url"`
	InstagramURL     string `gorm:"size:300" json:"instagram_url"`
	TiktokURL        string `gorm:"size:300" json:"tiktok_url"`
	StravaURL        string `gorm:"size:300" json:"strava_url"`
	WebsiteURL       string `gorm:"size:300" json:"website_url"`
	ResumeURL        string `gorm:"size:500" json:"resume_url"`
	YearsExperience  int    `gorm:"default:0" json:"years_experience"`
	AvailableForHire bool   `gorm:"default:true" json:"available_for_hire"`
	IsPublished      bool   `gorm:"default:true" json:"is_published"`
}

// ─── Tech Stack Category ───────────────────────────────────────────────────────

type StackCategory struct {
	Base
	Name      string      `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug      string      `gorm:"size:120;not null;uniqueIndex" json:"slug"`
	IconURL   string      `gorm:"size:300" json:"icon_url"`
	Color     string      `gorm:"size:20" json:"color"` // hex, e.g. "#3B82F6"
	SortOrder int         `gorm:"default:0" json:"sort_order"`
	Items     []StackItem `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// ─── Tech Stack Item ───────────────────────────────────────────────────────────

type StackItem struct {
	Base
	CategoryID uuid.UUID     `gorm:"type:uuid;not null;index" json:"category_id"`
	Category   StackCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name       string        `gorm:"size:100;not null" json:"name"`
	Slug       string        `gorm:"size:120;not null;uniqueIndex" json:"slug"`
	IconURL    string        `gorm:"size:300" json:"icon_url"`
	SortOrder  int           `gorm:"default:0" json:"sort_order"`
	IsPublished bool         `gorm:"default:true" json:"is_published"`
}

// ─── Project Category ──────────────────────────────────────────────────────────
type ProjectCategory struct {
	Base
	Name      string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug      string `gorm:"size:120;not null;uniqueIndex" json:"slug"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsPublished bool `gorm:"default:true" json:"is_published"`
}

// ─── Project ───────────────────────────────────────────────────────────────────

type Project struct {
	Base
	Title        string  `gorm:"size:200;not null" json:"title"`
	Slug         string  `gorm:"size:250;uniqueIndex;not null" json:"slug"`
	Description  string  `gorm:"type:text" json:"description"`
	Content      string  `gorm:"type:text" json:"content"` // rich text / markdown
	ThumbnailURL string  `gorm:"size:500" json:"thumbnail_url"`
	
	// Links
	DemoURL    string `gorm:"size:300" json:"demo_url"`
	RepoURL    string `gorm:"size:300" json:"repo_url"`
	DocURL     string `gorm:"size:300" json:"doc_url"`

	// Case study fields
	Problem  string `gorm:"type:text" json:"problem"`
	Solution string `gorm:"type:text" json:"solution"`
	MyRole   string `gorm:"type:text" json:"my_role"`
	Impact   string `gorm:"type:text" json:"impact"`

	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IsFeatured  bool       `gorm:"default:false" json:"is_featured"`
	IsPublished bool       `gorm:"default:true" json:"is_published"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`

	// Category (opsional)
	CategoryID *uuid.UUID       `gorm:"type:uuid;index" json:"category_id"` 
	Category   *ProjectCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// Relations
	Skills     []Skill     `gorm:"many2many:project_skills" json:"skills,omitempty"`
	StackItems []StackItem    `gorm:"many2many:project_stacks" json:"stack_items,omitempty"`
	Assets []Asset `gorm:"polymorphic:Owner;" json:"assets,omitempty"`
}

// ─── Project ↔ StackItem pivot ─────────────────────────────────────────────────

// GORM auto-creates this table via many2many tag above.
// Nama tabel: project_stacks
type AssetType string
const (
    AssetPhoto AssetType = "photo"
    AssetVideo AssetType = "video"
    AssetPDF   AssetType = "pdf"
    AssetDoc   AssetType = "doc"
)

type Asset struct {
    Base
    OwnerType string    `gorm:"size:50;not null;index" json:"owner_type"` // "project" / "experience"
    OwnerID   uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`
    Type      AssetType `gorm:"size:10;not null;default:'photo'" json:"type"`
    URL       string    `gorm:"size:500;not null" json:"url"`
    Title     string    `gorm:"size:200" json:"title"`     // untuk dokumen
    Caption   string    `gorm:"size:300" json:"caption"`  // untuk foto/video
    SortOrder int       `gorm:"default:0" json:"sort_order"`
}

// ─── Skill ─────────────────────────────────────────────────────────────────────

type Skill struct {
	Base
	Name        string `gorm:"size:100;not null" json:"name"`
	Category    string `gorm:"size:100" json:"category"` // legacy — tetap dipakai, Sprint 2 akan extend
	IconURL     string `gorm:"size:300" json:"icon_url"`
	Proficiency int    `gorm:"default:0" json:"proficiency"` // 0-100
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsPublished bool   `gorm:"default:true" json:"is_published"`
	Projects []Project `gorm:"many2many:project_skills" json:"projects,omitempty"`
}

// ─── Experience Category ───────────────────────────────────────────────────────
type ExperienceCategory struct {
	Base
	Name      string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug      string `gorm:"size:120;not null;uniqueIndex" json:"slug"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsPublished bool `gorm:"default:true" json:"is_published"`
}

// ─── Experience ────────────────────────────────────────────────────────────────

type Experience struct {
	Base
	Company     string     `gorm:"size:200;not null" json:"company"`
	Position    string     `gorm:"size:200;not null" json:"position"`
	Description string     `gorm:"type:text" json:"description"`
	LogoURL     string     `gorm:"size:300" json:"logo_url"`
	CompanyURL  string     `gorm:"size:300" json:"company_url"`
	Location    string     `gorm:"size:100" json:"location"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"` // nil = present
	IsCurrent   bool       `gorm:"default:false" json:"is_current"`
	IsPublished bool       `gorm:"default:true" json:"is_published"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	// Category (opsional)
	CategoryID *uuid.UUID          `gorm:"type:uuid;index" json:"category_id"`
	Category   *ExperienceCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	StackItems []StackItem `gorm:"many2many:experience_stacks" json:"stack_items,omitempty"`
	Skills    []Skill            `gorm:"many2many:experience_skills" json:"skills,omitempty"`
	Assets []Asset `gorm:"polymorphic:Owner;" json:"assets,omitempty"`
}

// ─── Education ─────────────────────────────────────────────────────────────────

type Education struct {
	Base
	Institution  string     `gorm:"size:200;not null" json:"institution"`
	Degree       string     `gorm:"size:200" json:"degree"`
	FieldOfStudy string     `gorm:"size:200" json:"field_of_study"`
	Description  string     `gorm:"type:text" json:"description"`
	LogoURL      string     `gorm:"size:300" json:"logo_url"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	GPA          string     `gorm:"size:10" json:"gpa"`
	IsPublished  bool       `gorm:"default:true" json:"is_published"`
	SortOrder    int        `gorm:"default:0" json:"sort_order"`
}

// ─── Contact (form submissions from public) ────────────────────────────────────

type Contact struct {
	Base
	Name    string `gorm:"size:100;not null" json:"name"`
	Email   string `gorm:"size:150;not null" json:"email"`
	Subject string `gorm:"size:250" json:"subject"`
	Message string `gorm:"type:text;not null" json:"message"`
	IsRead  bool   `gorm:"default:false" json:"is_read"`
}

// ─── Running Activity ──────────────────────────────────────────────────────────

type RunningActivity struct {
	Base
	Date           time.Time `gorm:"not null;index" json:"date"`
	Title          string    `gorm:"size:200" json:"title"`
	Notes          string    `gorm:"type:text" json:"notes"`
	MapImageURL    string    `gorm:"size:500" json:"map_image_url"`

	// Distance & time
	DurationSec    int     `gorm:"default:0" json:"duration_sec"`    // detik
	DistanceMeters float64 `gorm:"default:0" json:"distance_meters"` // meter

	// Calories
	TotalCalories  float64 `gorm:"default:0" json:"total_calories"`
	ActiveCalories float64 `gorm:"default:0" json:"active_calories"`

	// Pace & speed
	AvgPaceSec  int     `gorm:"default:0" json:"avg_pace_sec"`  // detik/km, tampilkan sebagai "5:30/km" di frontend
	AvgSpeedKph float64 `gorm:"default:0" json:"avg_speed_kph"` // km/h

	// Cadence & stride
	AvgCadence int     `gorm:"default:0" json:"avg_cadence"` // steps/min
	AvgStride  float64 `gorm:"default:0" json:"avg_stride"`  // meter
	Steps      int     `gorm:"default:0" json:"steps"`

	// Heart rate
	AvgHeartRate int `gorm:"default:0" json:"avg_heart_rate"` // bpm

	IsPublished bool `gorm:"default:true" json:"is_published"`
}